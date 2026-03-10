# Disable Template and YAML Function Resolution in Provider — Fix `!terraform.output` "text file busy" Crash

**Affected Versions:** `cloudposse/utils` provider v2.0.0–v2.0.2 (Atmos v1.207.0+ embedded)

**Severity:** Critical — provider crashes with "text file busy" when stack YAML contains
`!terraform.output` tags, breaking `terraform plan` / `terraform apply`

## Symptoms

Components using `remote-state` fail during plan when the referenced stack YAML files contain
`!terraform.output` YAML function calls:

```text
Error: failed to get terraform output for component vpc in stack dev-use1-frontoffice,
output vpc_default_security_group_id: failed to execute terraform output for component vpc
in stack dev-use1-frontoffice: terraform init failed: exit status 1

│ Error: Failed to install provider
│ Error while installing hashicorp/aws v5.100.0: open
│ /localhost/.terraform.d/plugin-cache/registry.opentofu.org/hashicorp/aws/5.100.0/linux_arm64/
│ terraform-provider-aws: text file busy
```

Occurs on Linux where the `ETXTBSY` errno is enforced. Does not reproduce
on macOS, which does not enforce this check.

## Root Cause

The provider v2.0.0 upgraded the embedded Atmos from v1.189.0 to v1.207.0. This newer Atmos
**resolves `!terraform.output` YAML tags** during `ProcessComponentInStack` by spawning child
`terraform init` + `terraform output` processes.

The resolution chain:

```text
processComponentInStackWithConfig()                          [component_processor.go:111]
  → e.ExecuteDescribeComponent(                              [hardcoded on lines 115-116]
        ProcessTemplates: true,
        ProcessYamlFunctions: true,
    )
    → Go template resolution (ProcessTemplates)              [Gomplate/Sprig/Atmos templates]
    → processTagTerraformOutput() (ProcessYamlFunctions)     [YAML tag parser]
      → outputGetter.GetOutput()                             [delegates to pkg/terraform/output]
        → runInit()                                          [terraform init via terraform-exec]
        → runOutput()                                        [terraform output]
```

Both `ProcessTemplates` and `ProcessYamlFunctions` are hardcoded to `true` in
`processComponentInStackWithConfig`. Template processing (`ProcessTemplates`) resolves
Go templates (Gomplate/Sprig/Atmos functions) in component configurations — but only when
templates are also enabled in `atmos.yaml` via `templates.settings.enabled: true`. YAML
function processing (`ProcessYamlFunctions`) resolves custom YAML tags like
`!terraform.output`, `!terraform.state`, and `!store`.

The child `terraform init` processes spawned by `!terraform.output` try to install providers
into the shared plugin cache (`TF_PLUGIN_CACHE_DIR` or `plugin_cache_dir`). On Linux, writing
to a binary file that is currently being executed by another process fails with `ETXTBSY`
("text file busy"). The outer OpenTofu process is already executing provider binaries from
this same cache directory.

In v1.189.0, neither templates nor YAML functions were resolved during
`ProcessComponentInStack` — tags were treated as opaque strings. The v2.0.0 upgrade introduced
eager resolution by hardcoding both `ProcessTemplates: true` and `ProcessYamlFunctions: true`
in `processComponentInStackWithConfig`.

### How Processing Is Controlled Across Entry Points

Template and YAML function processing can be controlled at multiple levels:

| Entry Point                                      | Templates                                                                                                                                                                                | YAML Functions                 | How to Disable                                                                                              |
|--------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------|-------------------------------------------------------------------------------------------------------------|
| **`atmos.yaml`**                                 | Requires `templates.settings.enabled: true` to activate. Sprig and Gomplate can be individually enabled via `templates.settings.sprig.enabled` and `templates.settings.gomplate.enabled` | Always active (no config gate) | Set `templates.settings.enabled: false`                                                                     |
| **Atmos CLI** (`atmos describe component`)       | Enabled by default                                                                                                                                                                       | Enabled by default             | `--process-templates=false` and `--process-functions=false` CLI flags                                       |
| **Stack imports**                                | Enabled by default                                                                                                                                                                       | N/A                            | `skip_templates_processing: true` on individual imports                                                     |
| **`ProcessComponentInStack` API** (this fix)     | Enabled by default (`true`)                                                                                                                                                              | Enabled by default (`true`)    | Pass `ProcessComponentInStackOptions` with `ProcessTemplates: &false` and/or `ProcessYamlFunctions: &false` |
| **`ProcessComponentFromContext` API** (this fix) | Enabled by default (`true`)                                                                                                                                                              | Enabled by default (`true`)    | Set `ProcessTemplates` and/or `ProcessYamlFunctions` fields on `ComponentFromContextParams`                 |
| **Provider** (after this fix)                    | **Disabled**                                                                                                                                                                             | **Disabled**                   | Hardcoded to `false` — provider does not need resolved values                                               |

## Why the Provider Does Not Need Template or YAML Function Resolution

The `utils_component_config` data source is used by the `remote-state` module to look up:

- Component backend configuration (S3 bucket, key pattern, workspace)
- Component workspace name
- Component variables (vars)

None of these require:

- **Template resolution** (`ProcessTemplates`) — Go template expressions (Gomplate/Sprig/Atmos
  functions) in component configs are meant for Atmos CLI rendering, not for the provider.
  The provider only needs the raw backend and workspace values, which are plain strings.
- **YAML function resolution** (`ProcessYamlFunctions`) — `!terraform.output`,
  `!terraform.state`, and `!store` tags return output values from other components, which are
  irrelevant for locating the remote state backend.

Disabling both avoids spawning child processes and unnecessary template evaluation inside
the provider plugin.

## Fix

### Atmos library change (`pkg/describe/component_processor.go`)

Add optional processing controls using a variadic options pattern for `ProcessComponentInStack`
and `*bool` fields on `ComponentFromContextParams`. This is **fully backward compatible** —
existing callers compile without changes, and omitted/nil values default to `true`.

**New `ProcessComponentInStackOptions` struct:**

```go
type ProcessComponentInStackOptions struct {
    ProcessTemplates     *bool // Controls Go template resolution. Defaults to true if nil.
    ProcessYamlFunctions *bool // Controls YAML function resolution. Defaults to true if nil.
}
```

**`ProcessComponentInStack`** — add variadic options (existing 4-arg calls still compile):

```go
// Before:
func ProcessComponentInStack(
    component string,
    stack string,
    atmosCliConfigPath string,
    atmosBasePath string,
) (map[string]any, error) {

// After (variadic — backward compatible):
func ProcessComponentInStack(
    component string,
    stack string,
    atmosCliConfigPath string,
    atmosBasePath string,
    opts ...ProcessComponentInStackOptions,
) (map[string]any, error) {
```

**`ComponentFromContextParams`** — add optional fields (existing struct literals still compile):

```go
type ComponentFromContextParams struct {
    Component            string
    Namespace            string
    Tenant               string
    Environment          string
    Stage                string
    AtmosCliConfigPath   string
    AtmosBasePath        string
    ProcessTemplates     *bool // Optional: defaults to true if nil
    ProcessYamlFunctions *bool // Optional: defaults to true if nil
}
```

**`processComponentInStackWithConfig`** — accept and use both parameters via `boolDefault`
helper (`nil` → `true`):

```go
func processComponentInStackWithConfig(
    atmosConfig *schema.AtmosConfiguration,
    component string,
    stack string,
    processTemplates *bool,
    processYamlFunctions *bool,
) (map[string]any, error) {
    return e.ExecuteDescribeComponent(&e.ExecuteDescribeComponentParams{
        AtmosConfig:          atmosConfig,
        Component:            component,
        Stack:                stack,
        ProcessTemplates:     boolDefault(processTemplates, true),
        ProcessYamlFunctions: boolDefault(processYamlFunctions, true),
    })
}
```

All existing callers (Atmos CLI, tests) continue to work without changes — the variadic
options are not passed, so both flags default to `true`.

### Provider change (`internal/provider/data_source_component_config.go`)

Pass `false` for both `processTemplates` and `processYamlFunctions` in both call sites:

```go
processTemplates := false
processYamlFunctions := false

atmosMu.Lock()
if len(stack) > 0 {
    result, err = p.ProcessComponentInStack(
        component, stack, atmosCliConfigPath, atmosBasePath,
        p.ProcessComponentInStackOptions{
            ProcessTemplates:     &processTemplates,
            ProcessYamlFunctions: &processYamlFunctions,
        },
    )
} else {
    result, err = p.ProcessComponentFromContext(&p.ComponentFromContextParams{
        Component:            component,
        Namespace:            namespace,
        Tenant:               tenant,
        Environment:          environment,
        Stage:                stage,
        AtmosCliConfigPath:   atmosCliConfigPath,
        AtmosBasePath:        atmosBasePath,
        ProcessTemplates:     &processTemplates,
        ProcessYamlFunctions: &processYamlFunctions,
    })
}
atmosMu.Unlock()
```

### Test changes

No changes needed for existing Atmos tests — the variadic options pattern is fully backward
compatible. New tests were added in the Atmos library to verify:

- `TestBoolDefault` — unit test for the `boolDefault` helper
- `TestProcessComponentInStackWithOptionsDisabled` — verifies `false` flags return results
  with vars/backend/workspace (no template or YAML function resolution)
- `TestProcessComponentInStackWithOptionsNilDefaultsToTrue` — verifies nil defaults to `true`
- `TestProcessComponentInStackBackwardCompatNoOptions` — verifies old 4-arg call still works
- `TestProcessComponentInStackDisabledMatchesEnabled` — verifies same vars with flags disabled
  (for configs without templates/YAML functions)
- `TestProcessComponentFromContextWithOptionsDisabled` — verifies struct-based API with both
  flags disabled

## Rollout

1. **Atmos**: Merge the `processYamlFunctions` parameter change, release new version (e.g. v1.209.0)
2. **Provider**: Update `go.mod` to new Atmos version, add `processYamlFunctions: false`, release v2.1.0
3. **Downstream**: Update `remote-state` module to pin `utils >= 2.1.0`

## Alternative Approaches Considered

| Approach                                                                                | Pros                                  | Cons                                                                        |
|-----------------------------------------------------------------------------------------|---------------------------------------|-----------------------------------------------------------------------------|
| **Disable plugin cache for child processes** (unset `TF_PLUGIN_CACHE_DIR` in child env) | Simple, targeted fix                  | Still spawns unnecessary child processes; slower; fragile                   |
| **Use temp cache for child processes** (`/tmp/atmos-utils-cache-XXXXX`)                 | Preserves caching                     | Still spawns unnecessary child processes; disk cleanup needed               |
| **Override `!terraform.output` handler to use `!terraform.state` path**                 | Transparent to users                  | Complex; couples provider to state reader internals                         |
| **Disable template and YAML function resolution entirely** (this fix)                   | No child processes; fastest; simplest | Template/YAML function values not available in provider output (not needed) |

The chosen approach (disable both template and YAML function resolution) is the most robust
because:

- The provider does not need resolved template or YAML function values
- It eliminates the entire class of child-process-related failures
- It avoids unnecessary template evaluation overhead inside the provider plugin
- It is the smallest code change with the clearest semantics
- It has no performance cost (skipping resolution is faster)

## References

- Hardcoded `ProcessTemplates: true` and `ProcessYamlFunctions: true`: Atmos `pkg/describe/component_processor.go:115-116`
- Provider call site: `internal/provider/data_source_component_config.go:119-130`
- `ExecuteDescribeComponentParams` struct: Atmos `internal/exec/describe_component.go:202-211`
- Concurrent `ReadDataSource` serialization: `internal/provider/atmos_lock.go`

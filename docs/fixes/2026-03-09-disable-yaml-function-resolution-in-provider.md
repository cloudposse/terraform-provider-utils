# Disable Template and YAML Function Resolution in Provider — Fix `!terraform.output` "text file busy" Crash

**Affected Versions:** `cloudposse/utils` provider v2.0.0–v2.0.2 (Atmos v1.207.0+ embedded)

**Severity:** Critical — provider crashes with "text file busy" when stack YAML contains
`!terraform.output` tags, breaking `terraform plan` / `terraform apply`

## Symptoms

Components using `remote-state` fail during plan when the referenced stack YAML files contain
`!terraform.output` YAML function calls:

```text
Error: failed to get terraform output for component vpc in stack dev-use1-frontend,
output vpc_default_security_group_id: failed to execute terraform output for component vpc
in stack dev-use1-frontend: terraform init failed: exit status 1

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

Add optional processing controls using a functional options pattern for both
`ProcessComponentInStack` and `ProcessComponentFromContext`. This is **fully backward
compatible** — existing callers compile without changes, and omitted options default to `true`.

**Functional options:**

```go
type ProcessOption func(*processOptions)

func WithProcessTemplates(enabled bool) ProcessOption
func WithProcessYamlFunctions(enabled bool) ProcessOption
```

**`ProcessComponentInStack`** — add variadic functional options (existing 4-arg calls still compile):

```go
// Before:
func ProcessComponentInStack(
    component string,
    stack string,
    atmosCliConfigPath string,
    atmosBasePath string,
) (map[string]any, error)

// After (variadic functional options — backward compatible):
func ProcessComponentInStack(
    component string,
    stack string,
    atmosCliConfigPath string,
    atmosBasePath string,
    opts ...ProcessOption,
) (map[string]any, error)
```

**`ProcessComponentFromContext`** — add variadic functional options:

```go
// Before:
func ProcessComponentFromContext(params *ComponentFromContextParams) (map[string]any, error)

// After (variadic functional options — backward compatible):
func ProcessComponentFromContext(
    params *ComponentFromContextParams,
    opts ...ProcessOption,
) (map[string]any, error)
```

All existing callers (Atmos CLI, tests) continue to work without changes — the variadic
options are not passed, so both flags default to `true`.

### Provider change (`internal/provider/data_source_component_config.go`)

Pass `WithProcessTemplates(false)` and `WithProcessYamlFunctions(false)` in both call sites:

```go
atmosMu.Lock()
if len(stack) > 0 {
    result, err = p.ProcessComponentInStack(
        component, stack, atmosCliConfigPath, atmosBasePath,
        p.WithProcessTemplates(false),
        p.WithProcessYamlFunctions(false),
    )
} else {
    result, err = p.ProcessComponentFromContext(&p.ComponentFromContextParams{
        Component:          component,
        Namespace:          namespace,
        Tenant:             tenant,
        Environment:        environment,
        Stage:              stage,
        AtmosCliConfigPath: atmosCliConfigPath,
        AtmosBasePath:      atmosBasePath,
    },
        p.WithProcessTemplates(false),
        p.WithProcessYamlFunctions(false),
    )
}
atmosMu.Unlock()
```

### Test changes

#### Atmos library tests (`pkg/describe/component_processor_test.go`)

Each flag is tested independently against its own fixture, proving the two flags are wired
independently and do not interfere with each other:

| Test | What It Verifies |
|------|------------------|
| `TestProcessComponentInStackTemplatesDisabledOnly` | `WithProcessTemplates(false)` preserves raw Go template strings while YAML functions remain enabled |
| `TestProcessComponentInStackTemplatesEnabledOnly` | `WithProcessTemplates(true)` resolves Go templates while YAML functions are disabled |
| `TestProcessComponentInStackYamlFunctionsDisabledOnly` | `WithProcessYamlFunctions(false)` preserves raw YAML function tags while templates remain enabled |
| `TestProcessComponentInStackYamlFunctionsEnabledOnly` | `WithProcessYamlFunctions(true)` resolves YAML function tags while templates are disabled |
| `TestProcessComponentInStackBackwardCompatNoOptions` | Old 4-arg call (no options) still works and returns correct vars |
| `TestProcessComponentFromContextWithProcessingDisabled` | `ProcessComponentFromContext` respects `WithProcessTemplates(false)` functional option |

#### Provider tests (`internal/component/component_processor_test.go`)

| Test | What It Verifies |
|------|------------------|
| `TestComponentProcessorWithProcessingDisabled` | `ProcessComponentInStack` with both flags `false` returns backend, workspace, vars |
| `TestComponentProcessorFromContextWithProcessingDisabled` | `ProcessComponentFromContext` with both flags `false` returns backend, workspace, vars |
| `TestComponentProcessorDisabledMatchesEnabled` | For configs without templates/YAML functions, results are identical with flags enabled or disabled |
| `TestComponentProcessorConsistencyWithProcessingDisabled` | Both API paths return the same results when processing is disabled |

## v1 Backward Compatibility Analysis

This fix will be released as a **v1 minor version** (v1.32.0) rather than a v2 release, because
all Cloud Posse `remote-state` modules and Terraform components pin the provider to `< 2`.
Releasing as v2 would require updating every downstream module.

### Version history and behavioral timeline

| Provider Version | Atmos Version | Templates Processed | YAML Functions Processed | Status |
|-----------------|---------------|--------------------|-----------------------|--------|
| **v1.31.0** | v1.189.0 | No | No | Last stable v1 release |
| **v2.0.0–v2.0.2** | v1.207.0 | Yes (hardcoded) | Yes (hardcoded) | Broken — `ETXTBSY` crash on Linux |
| **v1.32.0** (this fix) | v1.209.0 | No (`WithProcessTemplates(false)`) | No (`WithProcessYamlFunctions(false)`) | Restores v1.31.0 behavior |

### What changed between v1.31.0 and this branch

The commits between v1.31.0 and this branch are:

1. `43053f7` — Added Go linting (CI/tooling only, no behavioral change)
2. `762b3bd` — Use atmos instead of Makefile (build tooling only)
3. `d306476` — Fix release workflow (CI only)
4. `d668a0c` — Update Atmos to v1.207.0 (moved import from `pkg/component` to `pkg/describe`)
5. `3b54de6` — Serialize Atmos library calls with mutex (concurrent safety fix)
6. `01465c9` — Revert goreleaser config to v1 format (CI only)
7. `e035573` — Delete .goreleaser.yml (CI only)
8. `d397136`–`dd3547d` — Update Atmos to v1.208.0 then v1.209.0 (dependency update)
9. This PR — Disable template and YAML function processing

### Compatibility assessment

| Aspect | v1.31.0 | v1.32.0 (this fix) | Breaking? |
|--------|---------|---------------------|-----------|
| **Terraform provider schema** | All data sources unchanged | Identical schema | No |
| **Templates processed** | No (atmos v1.189.0 didn't process them) | No (`WithProcessTemplates(false)`) | No — same behavior |
| **YAML functions processed** | No (atmos v1.189.0 didn't process them) | No (`WithProcessYamlFunctions(false)`) | No — same behavior |
| **Concurrent access** | No serialization | Serialized with mutex | No — strictly safer |
| **Go version** | 1.22 | 1.26 | No — provider is a compiled binary |
| **Backend/workspace/vars output** | Works | Works (all tests pass) | No |
| **Import path** | `pkg/component` | `pkg/describe` | No — internal implementation detail |
| **`ProcessComponentFromContext` signature** | Positional args | Struct params | No — internal implementation detail |

### Why this is safe as a v1 minor release

1. **The provider is a compiled binary.** Users install it as a Terraform plugin — they never
   import Go packages from it. All Go-level API changes (import paths, function signatures,
   struct fields) are internal implementation details invisible to users.

2. **The Terraform provider schema is unchanged.** All data source schemas (`utils_component_config`,
   `utils_stack_config_yaml`, `utils_describe_stacks`, `utils_spacelift_stack_config`,
   `utils_deep_merge_yaml`, `utils_deep_merge_json`) have identical attributes, types, and
   defaults between v1.31.0 and this branch.

3. **Processing behavior matches v1.31.0.** In v1.31.0 (atmos v1.189.0), templates and YAML
   functions were not processed — they were treated as opaque strings. This fix explicitly
   disables both with `WithProcessTemplates(false)` and `WithProcessYamlFunctions(false)`,
   restoring the exact same behavior. The v2.0.0 regression (hardcoded `true`) is bypassed.

4. **All existing tests pass.** The full provider test suite (`internal/provider`,
   `internal/component`, `internal/describe`, `internal/merge`, `internal/spacelift`,
   `internal/stack`, `internal/convert`) passes without modification.

5. **The mutex serialization is additive.** The concurrent `ReadDataSource` fix (#523) only
   makes the provider safer under concurrent access — it cannot break existing single-threaded
   usage.

## Rollout

1. **Atmos**: ~~Merge the functional options change, release new version~~ Done — merged as
   [PR #2161](https://github.com/cloudposse/atmos/pull/2161), released as **v1.209.0**
2. **Provider**: Update `go.mod` to atmos v1.209.0, pass `WithProcessTemplates(false)` and
   `WithProcessYamlFunctions(false)`, release as **v1.32.0** (minor v1 bump)
3. **Downstream**: No changes needed — `remote-state` modules pinned to `< 2` will
   automatically pick up v1.32.0 on next `terraform init -upgrade`

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

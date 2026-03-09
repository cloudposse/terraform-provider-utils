# Disable YAML Function Resolution in Provider — Fix `!terraform.output` "text file busy" Crash

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

Occurs on Linux (Spacelift runners) where the `ETXTBSY` errno is enforced. Does not reproduce
on macOS, which does not enforce this check.

## Root Cause

The provider v2.0.0 upgraded the embedded Atmos from v1.189.0 to v1.207.0. This newer Atmos
**resolves `!terraform.output` YAML tags** during `ProcessComponentInStack` by spawning child
`terraform init` + `terraform output` processes.

The resolution chain:

```
processComponentInStackWithConfig()                          [component_processor.go:111]
  → e.ExecuteDescribeComponent(ProcessYamlFunctions: true)   [hardcoded on line 116]
    → processTagTerraformOutput()                            [YAML tag parser]
      → outputGetter.GetOutput()                             [delegates to pkg/terraform/output]
        → runInit()                                          [terraform init via terraform-exec]
        → runOutput()                                        [terraform output]
```

The child `terraform init` processes try to install providers into the shared plugin cache
(`TF_PLUGIN_CACHE_DIR` or `plugin_cache_dir`). On Linux, writing to a binary file that is
currently being executed by another process fails with `ETXTBSY` ("text file busy"). The outer
OpenTofu process is already executing provider binaries from this same cache directory.

In v1.189.0, `!terraform.output` tags were either not resolved or treated as opaque strings.
The v2.0.0 upgrade introduced eager resolution by hardcoding `ProcessYamlFunctions: true` in
`processComponentInStackWithConfig`.

## Why the Provider Does Not Need YAML Function Resolution

The `utils_component_config` data source is used by the `remote-state` module to look up:

- Component backend configuration (S3 bucket, key pattern, workspace)
- Component workspace name
- Component variables (vars)

None of these require resolving `!terraform.output` or `!terraform.state` YAML function values.
Those functions return output values from other components, which are irrelevant for the
provider's purpose of locating the remote state backend.

## Fix

### Atmos library change (`pkg/describe/component_processor.go`)

Add a `processYamlFunctions` parameter to the public API functions so callers can control
whether YAML functions are resolved. Use `*bool` (pointer) so `nil` defaults to `true` for
backward compatibility with existing callers (Atmos CLI, tests).

**`ProcessComponentInStack`** — add 5th parameter:

```go
// Before:
func ProcessComponentInStack(
component string,
stack string,
atmosCliConfigPath string,
atmosBasePath string,
) (map[string]any, error) {

// After:
func ProcessComponentInStack(
component string,
stack string,
atmosCliConfigPath string,
atmosBasePath string,
processYamlFunctions *bool,
) (map[string]any, error) {
```

**`ComponentFromContextParams`** — add field:

```go
type ComponentFromContextParams struct {
Component            string
Namespace            string
Tenant               string
Environment          string
Stage                string
AtmosCliConfigPath   string
AtmosBasePath        string
ProcessYamlFunctions *bool // Optional: defaults to true if nil
}
```

**`processComponentInStackWithConfig`** — accept and use the parameter:

```go
func processComponentInStackWithConfig(
atmosConfig *schema.AtmosConfiguration,
component string,
stack string,
processYamlFunctions *bool,
) (map[string]any, error) {
resolveYaml := true
if processYamlFunctions != nil {
resolveYaml = *processYamlFunctions
}

return e.ExecuteDescribeComponent(&e.ExecuteDescribeComponentParams{
AtmosConfig:          atmosConfig,
Component:            component,
Stack:                stack,
ProcessTemplates:     true,
ProcessYamlFunctions: resolveYaml,
})
}
```

All existing callers (Atmos CLI, tests) pass `nil` to preserve current behavior.

### Provider change (`internal/provider/data_source_component_config.go`)

Pass `false` for `processYamlFunctions` in both call sites:

```go
processYamlFunctions := false

atmosMu.Lock()
if len(stack) > 0 {
result, err = p.ProcessComponentInStack(
component, stack, atmosCliConfigPath, atmosBasePath,
&processYamlFunctions,
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
ProcessYamlFunctions: &processYamlFunctions,
})
}
atmosMu.Unlock()
```

### Test changes

Update `internal/component/component_processor_test.go` and the Atmos library's
`pkg/describe/component_processor_test.go` to pass `nil` as the new parameter in all
`ProcessComponentInStack` calls.

## Rollout

1. **Atmos**: Merge the `processYamlFunctions` parameter change, release new version (e.g. v1.209.0)
2. **Provider**: Update `go.mod` to new Atmos version, add `processYamlFunctions: false`, release v2.1.0
3. **Downstream**: Update `remote-state` module to pin `utils >= 2.1.0`

## Alternative Approaches Considered

| Approach                                                                                | Pros                                  | Cons                                                               |
|-----------------------------------------------------------------------------------------|---------------------------------------|--------------------------------------------------------------------|
| **Disable plugin cache for child processes** (unset `TF_PLUGIN_CACHE_DIR` in child env) | Simple, targeted fix                  | Still spawns unnecessary child processes; slower; fragile          |
| **Use temp cache for child processes** (`/tmp/atmos-utils-cache-XXXXX`)                 | Preserves caching                     | Still spawns unnecessary child processes; disk cleanup needed      |
| **Override `!terraform.output` handler to use `!terraform.state` path**                 | Transparent to users                  | Complex; couples provider to state reader internals                |
| **Disable YAML function resolution entirely** (this fix)                                | No child processes; fastest; simplest | YAML function values not available in provider output (not needed) |

The chosen approach (disable YAML function resolution) is the most robust because:

- The provider does not need resolved YAML function values
- It eliminates the entire class of child-process-related failures
- It is the smallest code change with the clearest semantics
- It has no performance cost (skipping resolution is faster)

## References

- Issue analysis: `docs/fixes/dwt-stacks-opentofu-utils-provider-issues.md`
- Hardcoded `ProcessYamlFunctions: true`: Atmos `pkg/describe/component_processor.go:116`
- Provider call site: `internal/provider/data_source_component_config.go:119-130`
- `ExecuteDescribeComponentParams` struct: Atmos `internal/exec/describe_component.go:202-211`
- Concurrent `ReadDataSource` serialization: `internal/provider/atmos_lock.go`

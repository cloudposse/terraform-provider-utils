# Atmos Version Mismatch Causes `utils_component_config` Plugin Crash

**Reported by:** Community user

**Affected Versions:** `cloudposse/utils` provider v1.31.0 (embeds Atmos v1.189.0)

**Severity:** Critical - `data "utils_component_config"` crashes with "Plugin did not respond", blocking all
components that use the `cloudposse/stack-config/yaml//modules/remote-state` module

## Symptoms

Running `atmos terraform plan aws-sso -s core-gbl-root` (or any component that uses remote-state modules)
fails with plugin crashes:

```text
module.role_map.module.account_map.data.utils_component_config.config[0]: Still reading... [10s elapsed]
module.tfstate.data.utils_component_config.config[0]: Still reading... [10s elapsed]
module.iam_roles.module.account_map.data.utils_component_config.config[0]: Still reading... [10s elapsed]

Planning failed. Terraform encountered an error while generating this plan.

Error: Plugin did not respond

  with module.iam_roles.module.account_map.data.utils_component_config.config[0],
  on .terraform/modules/iam_roles.account_map/modules/remote-state/main.tf line 1,
  in data "utils_component_config" "config":
   1: data "utils_component_config" "config" {

The plugin encountered an error, and failed to respond to the
plugin.(*GRPCProvider).ReadDataSource call. The plugin logs may contain more details.
```

Notably, some parallel `utils_component_config` reads succeed (2 out of 5 in the reported case),
while the remaining 3 crash after 10+ seconds. The "Plugin did not respond" error means the provider
process panicked or exited unexpectedly.

## Root Cause

The `cloudposse/utils` Terraform provider v1.31.0 has Atmos **v1.189.0** compiled into its binary
(`go.mod`: `github.com/cloudposse/atmos v1.189.0`). When any `data "utils_component_config"`,
`data "utils_describe_stacks"`, or `data "utils_stack_config_yaml"` executes, the provider calls
`cfg.InitCliConfig()` which parses the project's `atmos.yaml` and processes all stack configurations.

The user's infrastructure uses Atmos CLI v1.200+ and their `atmos.yaml` and stack
files contain features that did not exist in Atmos v1.189.0:

- **`stores`** block (AWS SSM Parameter Store configuration) - new top-level config section
- **`hooks`** with `store-outputs` command - new lifecycle hook feature
- **`templates.settings.gomplate`** - new template engine configuration
- **`!terraform.state` / `!terraform.output`** YAML custom tags with Go template expressions

When the old v1.189.0 schema structs encounter these unknown configuration sections during YAML
unmarshaling, the provider panics (nil pointer dereference or unexpected type assertion), killing the
gRPC plugin process.

### Why some calls succeed and others fail

Terraform executes data source reads concurrently. The `aws-sso` component triggers 5 parallel
`utils_component_config` calls (via `module.account_map`, `module.iam_roles.account_map`,
`module.iam_roles_root.account_map`, `module.role_map.account_map`, `module.tfstate`). The first two
calls can complete before hitting the panic path, while the remaining three crash during concurrent
stack processing with the incompatible configuration.

### Call chain

1. Terraform invokes `ReadDataSource` on the `cloudposse/utils` provider plugin (gRPC)
2. `dataSourceComponentConfigRead()` in `internal/provider/data_source_component_config.go`
3. Calls `p.ProcessComponentInStack(component, stack, ...)` from `atmos/pkg/component`
4. Internally calls `cfg.InitCliConfig(...)` which parses `atmos.yaml`
5. The v1.189.0 config parser encounters unknown sections (`stores`, `hooks`, `templates.settings.gomplate`)
6. Panic -> gRPC process exits -> Terraform reports "Plugin did not respond"

## Scope of Impact

### Data sources affected

All 7 data sources in the provider are at risk because they all call into Atmos packages that
parse `atmos.yaml`:

| Data Source                       | Atmos API Used                                          | Status in v1.206.2          |
|-----------------------------------|---------------------------------------------------------|-----------------------------|
| `utils_component_config`          | `pkg/component.ProcessComponentInStack`                 | **API deleted in v1.201.0** |
| `utils_component_config`          | `pkg/component.ProcessComponentFromContext`             | **API deleted in v1.201.0** |
| `utils_describe_stacks`           | `pkg/describe.ExecuteDescribeStacks`                    | Exists, signature unchanged |
| `utils_describe_stacks`           | `pkg/config.InitCliConfig`                              | Exists, signature unchanged |
| `utils_describe_stacks`           | `pkg/config.GetStackNameFromContextAndStackNamePattern` | Exists                      |
| `utils_stack_config_yaml`         | `pkg/stack.ProcessYAMLConfigFiles`                      | Exists, signature unchanged |
| `utils_stack_config_yaml`         | `pkg/config.InitCliConfig`                              | Exists, signature unchanged |
| `utils_spacelift_stack_config`    | `pkg/spacelift.CreateSpaceliftStacks`                   | Exists, signature unchanged |
| `utils_aws_eks_update_kubeconfig` | `pkg/aws.ExecuteAwsEksUpdateKubeconfig`                 | Exists                      |
| `utils_deep_merge_json`           | `pkg/merge.MergeWithOptions`                            | Exists                      |
| `utils_deep_merge_yaml`           | `pkg/merge.MergeWithOptions`                            | Exists                      |

### Upgrade blocker: deleted API in `pkg/component`

In Atmos v1.201.0 (commit `a10d9ad23`, PR #1774 — "Path-based component resolution for all
commands"), the file `pkg/component/component_processor.go` was **deleted entirely**. This file
contained the two public functions the provider depends on:

- `ProcessComponentInStack(component, stack, atmosCliConfigPath, atmosBasePath string) (map[string]any, error)`
-
`ProcessComponentFromContext(component, namespace, tenant, environment, stage, atmosCliConfigPath, atmosBasePath string) (map[string]any, error)`

The `pkg/component` package was redesigned into a component registry/provider abstraction:

- `provider.go` - `ComponentProvider` interface + `ExecutionContext` struct
- `registry.go` - thread-safe `ComponentRegistry`
- `resolver.go` - path-based component resolution

The original processing logic moved to `internal/exec.ProcessStacks()`, which is an **internal
package** and cannot be imported by external Go modules.

### User's affected components

The crash affects any component using
`cloudposse/stack-config/yaml//modules/remote-state` (version 1.8.0), which internally calls
`data "utils_component_config"`. This module is used by **56 components** across the codebase,
including `aws-sso`, `vpc`, `ecs`, `aurora-postgres`, `dns-delegated`, and many more.

Two components also declare the `cloudposse/utils` provider directly:

- `account-map` (version `>= 1.10.0`) - uses `data "utils_describe_stacks"`
- `spacelift/admin-stack` (version `>= 1.14.0`) - uses utils via the `spacelift-stacks-from-atmos-config` module

## Fix

The fix spans two repos: Atmos (restore public API) and this provider (update import + bump version).

### 1. Atmos: restore public API in `pkg/describe`

The functions cannot be restored in `pkg/component` because `internal/exec` now imports `pkg/component`
(for the `ComponentProvider` interface), which would create an import cycle. Instead,
`ProcessComponentInStack` and `ProcessComponentFromContext` are added to `pkg/describe/component_processor.go`
as thin wrappers delegating to `internal/exec.ExecuteDescribeComponent`.

See Atmos doc: `docs/fixes/2026-02-15-restore-component-processor-public-api.md`

Atmos branch: `aknysh/update-for-utils-provider-1`

### 2. Provider: update import and bump Atmos version

**`internal/provider/data_source_component_config.go`** — change import:

```go
// Before:
p "github.com/cloudposse/atmos/pkg/component"

// After:
p "github.com/cloudposse/atmos/pkg/describe"
```

The function signatures are identical, so no other code changes are needed.

**`go.mod`** — bump Atmos dependency:

```
github.com/cloudposse/atmos v1.189.0 -> v1.206.3
```

The provider will compile once Atmos v1.206.3 is released with the restored API.

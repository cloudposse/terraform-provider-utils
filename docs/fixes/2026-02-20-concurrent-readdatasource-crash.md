# Concurrent `ReadDataSource` Crashes Provider — `os.Exit(1)` from Thread-Unsafe Atmos Library

**Affected Versions:** `cloudposse/utils` provider v2.0.0 (Atmos v1.207.0 embedded)

**Severity:** Critical — provider process exits with code 1 during concurrent `ReadDataSource`
calls, producing "Plugin did not respond" errors in Terraform

## Symptoms

Components with multiple `data "utils_component_config"` data sources (e.g., 5 concurrent reads)
crash the provider during `terraform plan`:

```text
Error: Plugin did not respond

  with module.iam_roles.module.account_map.data.utils_component_config.config[0],
  on .terraform/modules/iam_roles.account_map/modules/remote-state/main.tf line 1,
  in data "utils_component_config" "config":
   1: data "utils_component_config" "config" {
```

`TF_LOG=TRACE` reveals the provider exits with **exit status 1** (not a panic or signal):

```text
provider: plugin process exited: path=...terraform-provider-utils pid=3934 error="exit status 1"
```

Components with 1–3 concurrent reads typically work. Components with 5+ concurrent reads crash
near-deterministically. The crash happens ~700ms after the first `ReadDataSource` request — all
concurrent reads fail, none returns a response.

## Root Cause

The Atmos library was designed as a single-threaded CLI tool. It has **package-level mutable
state** that is explicitly documented as not thread-safe:

```go
// pkg/config/load.go:51-54
// NOTE: This package-level state assumes sequential (non-concurrent) calls to LoadConfig.
// LoadConfig is NOT safe for concurrent use.
var mergedConfigFiles []string
```

Additional thread-unsafe global state includes:

- `errors/error_funcs.go:31` — `var atmosConfig *schema.AtmosConfiguration`
- `errors/error_funcs.go:28` — `var render *markdown.Renderer`
- `errors/error_funcs.go:34` — `var verboseFlag`

### Crash sequence

1. Terraform invokes multiple `ReadDataSource` calls concurrently (one per data source instance)
2. Each call enters `ProcessComponentInStack` → `InitCliConfig` → `LoadConfig`
3. `LoadConfig` resets and writes to the shared `mergedConfigFiles` slice (`mergedConfigFiles = nil`)
4. Concurrent goroutines corrupt the slice through interleaved reads/writes
5. Corrupted config causes errors during template or YAML function processing
6. Errors hit `CheckErrorPrintAndExit()` → `Exit()` → `os.Exit(1)`
7. `os.Exit(1)` kills the gRPC plugin process without returning an error to Terraform
8. Terraform reports "Plugin did not respond"

### Why `os.Exit` instead of error propagation

The Atmos library uses `CheckErrorPrintAndExit()` in many code paths within `internal/exec`:

- `utils.go:753,765` — template processing errors
- `yaml_func_store.go:29,89,99,107` — `!store` YAML tag errors
- `yaml_func_store_get.go:52,90,101,116` — `!store.get` YAML tag errors
- `describe_stacks.go:489,743,982` — describe stacks errors

This function is designed for CLI usage — it prints an error and exits the process. Inside a
Terraform provider (gRPC plugin), calling `os.Exit(1)` terminates the plugin without returning
a diagnostic error to Terraform.

### Debug log timeline (provider v2.0.0, pid=3934)

| Timestamp      | Event                                                        |
|----------------|--------------------------------------------------------------|
| `19:24:15.825` | Provider starts, configures mTLS                             |
| `19:24:15.941` | GetProviderSchema — success                                  |
| `19:24:15.945` | Configure — success                                          |
| `19:24:16.117` | ValidateDataSourceConfig — success                           |
| `19:24:16.122` | **ReadDataSource #1** — "Calling downstream" (never returns) |
| `19:24:16.152` | **ReadDataSource #2** — "Calling downstream" (never returns) |
| `19:24:16.276` | **ReadDataSource #3** — "Calling downstream" (never returns) |
| `19:24:16.817` | **Plugin process exited: exit status 1**                     |
| `19:24:16.817` | gRPC: "connection reset by peer"                             |

GetProviderSchema, Configure, and ValidateDataSourceConfig all succeed because they don't call
`LoadConfig`. The crash happens specifically during `ReadDataSource` → `ProcessComponentInStack`
→ `InitCliConfig` → `LoadConfig`.

## Fix

### Provider-side fix (this repo)

Add a package-level `sync.Mutex` in the provider to serialize all calls into the Atmos library.
This prevents concurrent goroutines from accessing the thread-unsafe global state simultaneously.

**New file: `internal/provider/atmos_lock.go`**

```go
package provider

import "sync"

// atmosMu serializes all calls into the Atmos library.
// The Atmos library uses package-level mutable state (e.g., mergedConfigFiles in
// pkg/config/load.go) that is explicitly documented as not safe for concurrent use.
// Terraform invokes ReadDataSource concurrently for independent data sources, so
// without this mutex, concurrent calls corrupt shared state and trigger os.Exit(1)
// via CheckErrorPrintAndExit, killing the gRPC plugin process.
var atmosMu sync.Mutex
```

**Modified data source files:**

Each `ReadContext` function wraps its Atmos library calls with `atmosMu.Lock()` /
`atmosMu.Unlock()`:

- `data_source_component_config.go` — wraps `ProcessComponentInStack` / `ProcessComponentFromContext`
- `data_source_describe_stacks.go` — wraps `InitCliConfig` + `ExecuteDescribeStacks`
- `data_source_stack_config_yaml.go` — wraps `InitCliConfig` + `ProcessYAMLConfigFiles`
- `data_source_spacelift_stack_config.go` — wraps `CreateSpaceliftStacks`

### Long-term fix (Atmos library)

The Atmos library should be refactored to:

1. Eliminate package-level mutable state in `pkg/config` and `errors`
2. Pass configuration through context or options structs instead of global variables
3. Replace `CheckErrorPrintAndExit` / `os.Exit` calls in library code paths with proper error
   returns, so embedded consumers (like this provider) can handle errors gracefully

## References

- Atmos `LoadConfig` thread-safety comment: `pkg/config/load.go:51-54`
- `CheckErrorPrintAndExit` implementation: `errors/error_funcs.go:324-366`
- Previous fix (version mismatch): `docs/fixes/2026-02-19-atmos-version-mismatch-plugin-crash.md`

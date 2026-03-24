# Base Path Resolution Fallback When Git Root Is Unavailable

**Atmos Version:** v1.211.0 (upgraded from v1.210.1)
**PR:** [cloudposse/atmos#2236](https://github.com/cloudposse/atmos/pull/2236)
**Release:** [v1.211.0](https://github.com/cloudposse/atmos/releases/tag/v1.211.0)
**Severity:** Medium — `failed to find import` errors when `ATMOS_BASE_PATH` is set to a
relative path on CI workers (e.g., Spacelift) that lack a `.git` directory

## Summary

Atmos v1.211.0 fixes base path resolution when `ATMOS_BASE_PATH` (or `--base-path` /
`atmos_base_path` provider parameter) is set to a relative path **and no `.git` directory
exists**. This is a follow-up to the v1.210.1
fix ([2026-03-18 fix](./2026-03-18-base-path-resolution-for-relative-paths.md)),
which added `os.Stat` + CWD fallback to `tryResolveWithGitRoot` but missed the same fix
in the `tryResolveWithConfigPath` fallback path.

## Root Cause

The v1.210.1 fix added `os.Stat` validation and CWD fallback to `tryResolveWithGitRoot()`.
However, when `getGitRootOrEmpty()` returns `""` (no `.git` directory — common on Spacelift
and other CI workers), the code falls through to `tryResolveWithConfigPath()`. That function
lacked the same `os.Stat` validation and CWD fallback — it unconditionally joined the path
with `cliConfigPath` (the `atmos.yaml` directory), producing a wrong/nonexistent path.

### Failure Path (Before Fix)

```
ATMOS_BASE_PATH=.terraform/modules/monorepo  (no .git directory)
  → tryResolveWithGitRoot("...", cliConfigPath)
  → getGitRootOrEmpty() returns ""
  → falls through to tryResolveWithConfigPath("...", cliConfigPath)
  → unconditionally returns filepath.Join(cliConfigPath, path)
  → wrong path → "failed to find import"
```

## What Changed in Atmos v1.211.0

### Modified Functions (All Unexported)

| Function                                                | Change                                                                      |
|---------------------------------------------------------|-----------------------------------------------------------------------------|
| `tryResolveWithGitRoot(path, cliConfigPath, source)`    | Added `source` parameter for source-aware resolution                        |
| `tryResolveWithConfigPath(path, cliConfigPath, source)` | Added `source` parameter + `os.Stat` validation + source-aware CWD fallback |
| `tryCWDRelative(path)` *(new)*                          | Extracted helper — checks if path exists relative to CWD                    |

### Source-Aware Fallback Ordering

When git root is unavailable, the resolution order now depends on path source:

| Source                                          | Fallback Order             | Rationale                                                        |
|-------------------------------------------------|----------------------------|------------------------------------------------------------------|
| **Runtime** (env var, CLI flag, provider param) | CWD first, then config dir | Shell convention — runtime paths are relative to where you run   |
| **Config** (`base_path` in `atmos.yaml`)        | Config dir first, then CWD | Config convention — config paths are relative to the config file |

### New Constants

- `basePathSourceRuntime = "runtime"` — replaces string literal
- `cwdResolutionErrFmt` — formatted error message for resolution failures

### Resolution Path (After Fix)

```
ATMOS_BASE_PATH=.terraform/modules/monorepo  (no .git directory, runtime source)
  → tryResolveWithGitRoot("...", cliConfigPath, "runtime")
  → getGitRootOrEmpty() returns ""
  → falls through to tryResolveWithConfigPath("...", cliConfigPath, "runtime")
  → source == "runtime" → tryCWDRelative(path) first
  → os.Stat(CWD/.terraform/modules/monorepo) → exists → returns CWD path
```

## Impact on terraform-provider-utils

### No Breaking Changes

1. **All modified functions are unexported** — `tryResolveWithGitRoot`, `tryResolveWithConfigPath`,
   and `tryCWDRelative` are internal to `pkg/config`. The provider never calls them directly.

2. **Public API is unchanged** — `InitCliConfig`, `ProcessComponentInStack`,
   `ProcessComponentFromContext`, `ExecuteDescribeStacks`, `MergeWithOptions` — all signatures
   are identical to v1.210.1.

3. **The fix is transparent** — provider users on CI workers without `.git` who set
   `ATMOS_BASE_PATH` to a relative path will stop getting `failed to find import` errors
   without any provider code changes.

### Provider Call Sites (Unchanged)

| Data Source              | Call                                                                                               | File                              |
|--------------------------|----------------------------------------------------------------------------------------------------|-----------------------------------|
| `utils_component_config` | `ProcessComponentInStack(component, stack, atmosCliConfigPath, atmosBasePath, ...)`                | `data_source_component_config.go` |
| `utils_component_config` | `ProcessComponentFromContext(&ComponentFromContextParams{AtmosBasePath: atmosBasePath, ...}, ...)` | `data_source_component_config.go` |
| `utils_describe_stacks`  | `InitCliConfig(ConfigAndStacksInfo{AtmosBasePath: atmosBasePath, ...}, true)`                      | `data_source_describe_stacks.go`  |

### Why This Is Safe for Existing Provider Users

1. **Empty `atmos_base_path` (default)** — the most common case — is completely unaffected.
   Empty paths return early in `tryResolveWithGitRoot` before any new code is reached.

2. **Config-file relative paths** (`../../examples/tests`) — classified as dot-prefixed,
   routed to `resolveDotPrefixPath`. Unchanged behavior.

3. **Users with `.git` directory present** — `getGitRootOrEmpty()` returns a valid path,
   so `tryResolveWithGitRoot` handles resolution before `tryResolveWithConfigPath` is called.
   Unchanged behavior.

4. **The new `os.Stat` + source-aware fallback in `tryResolveWithConfigPath`** is strictly
   additive — it only activates when git root is unavailable AND the old unconditional join
   would have produced an incorrect path.

### Risk Assessment

| Scenario                                 | Before (v1.210.1)                    | After (v1.211.0)                     | Breaking?  |
|------------------------------------------|--------------------------------------|--------------------------------------|------------|
| Empty base path, `.git` present          | Git root                             | Same                                 | No         |
| Empty base path, no `.git`               | Config dir -> CWD                    | Same                                 | No         |
| Relative path, `.git` present            | Git root + `os.Stat` fallback to CWD | Same                                 | No         |
| Relative path, no `.git`, runtime source | Config dir join (often wrong)        | **CWD first, then config dir**       | No (fix)   |
| Relative path, no `.git`, config source  | Config dir join                      | **Config dir + `os.Stat`, then CWD** | No (safer) |
| Absolute path                            | Pass through                         | Same                                 | No         |
| Dot-prefixed path                        | Source-aware resolution              | Same                                 | No         |

## Other Changes in v1.211.0

The release includes additional PRs that do not affect the provider's public API:

| PR                                                     | Description                                                                         | Provider Impact                          |
|--------------------------------------------------------|-------------------------------------------------------------------------------------|------------------------------------------|
| [#2225](https://github.com/cloudposse/atmos/pull/2225) | Refactored `processArgsAndFlags` — fixes for boolean flags, `SplitN`, strategy leak | None — `internal/exec` only              |
| [#2204](https://github.com/cloudposse/atmos/pull/2204) | Reduced `ExecuteDescribeStacks` cyclomatic complexity from 247 to 10                | Low — internal refactor, same public API |
| [#2226](https://github.com/cloudposse/atmos/pull/2226) | Refactored `ExecuteTerraform` complexity 160 to 9                                   | None — CLI execution only                |
| [#2243](https://github.com/cloudposse/atmos/pull/2243) | Better error messages for terraform component load errors                           | Positive — improved diagnostics          |
| [#2149](https://github.com/cloudposse/atmos/pull/2149) | EKS kubeconfig authentication integration                                           | None — new feature                       |
| [#2229](https://github.com/cloudposse/atmos/pull/2229) | Isolated browser sessions for `atmos auth console`                                  | None — CLI only                          |

## Tests

### Existing Tests (Must Pass)

All existing tests in the provider must continue to pass:

| Test File                                              | Tests                                                             |
|--------------------------------------------------------|-------------------------------------------------------------------|
| `internal/component/component_processor_test.go`       | Component processing, consistency, processing options, base paths |
| `internal/describe/describe_stacks_test.go`            | Stack description with filters                                    |
| `internal/spacelift/spacelift_stack_processor_test.go` | Spacelift config                                                  |
| `internal/stack/stack_processor_test.go`               | Stack processing                                                  |
| `internal/merge/merge_test.go`                         | Deep merge strategies                                             |
| `internal/convert/*_test.go`                           | YAML/JSON conversion                                              |
| `internal/provider/provider_test.go`                   | Provider schema                                                   |
| `internal/provider/provider_utils_test.go`             | Provider utilities                                                |

### New Tests in Atmos v1.211.0

| Test                                                 | What It Verifies                                                   |
|------------------------------------------------------|--------------------------------------------------------------------|
| `TestBasePathResolutionWithoutGitRoot_RuntimeSource` | Runtime-source relative path resolves via CWD when no `.git`       |
| `TestBasePathResolutionWithoutGitRoot_ConfigSource`  | Config-source relative path resolves via config dir when no `.git` |
| `TestBasePathResolutionWithoutGitRoot_BothExist`     | Correct precedence when path exists at both CWD and config dir     |
| `TestBasePathResolutionWithoutGitRoot_NeitherExist`  | Returns config-dir path for consistent error messages              |
| `TestBasePathResolutionWithoutGitRoot_AbsolutePath`  | Absolute paths pass through unchanged regardless of `.git`         |

### Cycle Detection Tests (Also Added in v1.211.0)

4 new tests in `internal/exec/stack_processor_utils_test.go` verify `metadata.component`
stack overflow detection — not directly related to the base path fix but included in this release.

## References

- Atmos PR: [cloudposse/atmos#2236](https://github.com/cloudposse/atmos/pull/2236)
- Atmos Release: [v1.211.0](https://github.com/cloudposse/atmos/releases/tag/v1.211.0)
- Previous fix: [2026-03-18-base-path-resolution-for-relative-paths.md](./2026-03-18-base-path-resolution-for-relative-paths.md)
- Previous Atmos PR: [cloudposse/atmos#2215](https://github.com/cloudposse/atmos/pull/2215)

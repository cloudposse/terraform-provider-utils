# Base Path Resolution for ATMOS_BASE_PATH and --base-path with Relative Paths

**Atmos Version:** v1.210.1 (upgraded from v1.209.0)
**PR:** [cloudposse/atmos#2215](https://github.com/cloudposse/atmos/pull/2215)
**Severity:** Medium — `failed to find import` errors when `ATMOS_BASE_PATH` or `atmos_base_path`
provider parameter is set to a relative path like `.terraform/modules/monorepo`

## Summary

Atmos v1.210.1 fixes base path resolution when `ATMOS_BASE_PATH` env var (or `--base-path`
flag / `atmos_base_path` provider parameter) is set to a relative path. Previously,
`resolveAbsolutePath()` routed simple relative paths through git root discovery, which
could produce incorrect paths when the CWD differed from the git root.

## What Changed in Atmos v1.210.1

### 4-Category Path Classification

Every `base_path` value is now classified into one of four categories:

| Category     | Pattern                                  | Resolution                                     |
|--------------|------------------------------------------|------------------------------------------------|
| **Empty**    | `""`, unset                              | Git root -> config dir -> CWD (smart default)  |
| **Dot**      | `"."`, `"./foo"`, `".."`, `"../foo"`     | Source-dependent anchor (see below)            |
| **Bare**     | `"foo"`, `"foo/bar"`, `".terraform/..."` | Git root search with `os.Stat` fallback to CWD |
| **Absolute** | `"/abs/path"`                            | Pass through unchanged                         |

### Source-Aware Resolution (`BasePathSource`)

A `BasePathSource` field on `AtmosConfiguration` tracks where the base path came from:

- **Runtime sources** (env var `ATMOS_BASE_PATH`, CLI flag `--base-path`, provider param
  `atmos_base_path`): set `BasePathSource = "runtime"`. Dot-prefixed paths resolve relative
  to **CWD** (shell convention).
- **Config source** (`base_path` in `atmos.yaml`): dot-prefixed paths resolve relative to
  the **directory containing `atmos.yaml`** (config-file convention).

### Git Root Fallback with `os.Stat` Validation

`tryResolveWithGitRoot()` now validates the git-root-joined path with `os.Stat`. If the
path doesn't exist at git root but does exist relative to CWD, it falls back to the
CWD-relative path.

### Actionable Error Messages

`failed to find import` errors now include context about which import failed and what path
was searched, with hints about checking `base_path`, `stacks.base_path`, and `ATMOS_BASE_PATH`.

## Impact on terraform-provider-utils

### Provider Call Sites

The provider has two data sources that pass `atmosBasePath`:

| Data Source              | Call                                                                                               | File                                      |
|--------------------------|----------------------------------------------------------------------------------------------------|-------------------------------------------|
| `utils_component_config` | `ProcessComponentInStack(component, stack, atmosCliConfigPath, atmosBasePath, ...)`                | `data_source_component_config.go:119`     |
| `utils_component_config` | `ProcessComponentFromContext(&ComponentFromContextParams{AtmosBasePath: atmosBasePath, ...}, ...)` | `data_source_component_config.go:124-135` |
| `utils_describe_stacks`  | `cfg.InitCliConfig(ConfigAndStacksInfo{AtmosBasePath: atmosBasePath, ...}, true)`                  | `data_source_describe_stacks.go:149-158`  |

### How the Provider Passes Base Path

The provider passes `atmos_base_path` to `configAndStacksInfo.AtmosBasePath`. In
`InitCliConfig` (config.go:38-41), if `AtmosBasePath != ""`, it sets
`atmosConfig.BasePath = configAndStacksInfo.AtmosBasePath` and
`atmosConfig.BasePathSource = "runtime"`.

### Internal Code Path Analysis

#### `resolveAbsolutePath` (config.go:227-253)

Three-way routing based on path classification:

1. **Absolute** (`filepath.IsAbs(path)`) — return as-is
2. **Dot-prefixed** (`"."`, `".."`, `"./..."`, `"../..."`) — routes to
   `resolveDotPrefixPath(path, cliConfigPath, source)` which is source-aware
3. **Everything else** (empty `""` or bare like `"stacks"`, `".terraform/..."`) — routes to
   `tryResolveWithGitRoot(path, cliConfigPath)`

Note: `.terraform/modules/monorepo` does NOT match dot-prefixed (`.t` != `./`). It is
classified as a bare path. `.terraform` is a directory name, not a relative path prefix.

#### `resolveDotPrefixPath` (config.go:269-283)

- `source == "runtime"` → `filepath.Abs(path)` (CWD-relative, shell convention)
- `source != "runtime"` + `cliConfigPath != ""` → `filepath.Abs(filepath.Join(cliConfigPath, path))` (
  config-dir-relative)
- `source != "runtime"` + `cliConfigPath == ""` → `filepath.Abs(path)` (CWD fallback)

#### `tryResolveWithGitRoot` (config.go:289-333)

1. Get git root. If empty → `tryResolveWithConfigPath` fallback
2. If `path == ""` → return git root directly (early return at line 297)
3. `os.Stat(gitRoot/path)` — if exists, return it
4. If `os.IsNotExist` → try `os.Stat(CWD/path)` — if exists, return it (NEW fallback)
5. Neither exists → return `gitRoot/path` (original behavior for consistent error messages)

#### Override Ordering (Precedence)

`processEnvVars` sets `BasePathSource = "runtime"` for `ATMOS_BASE_PATH` env var →
`setBasePaths` overrides for `--base-path` CLI flag →
`InitCliConfig` lines 38-41 override for `AtmosBasePath` struct field (provider param).
Last one wins. Correct precedence: provider param > CLI flag > env var > config file.

### Scenario-by-Scenario Analysis

#### Scenario A: Empty `atmos_base_path` (most common — default provider usage)

1. `InitCliConfig` line 38: `AtmosBasePath == ""` — condition false, skip
2. `BasePathSource` stays `""` (config source)
3. `BasePath` comes from `atmos.yaml` (e.g., `base_path: "../../examples/tests"` in test fixtures)
4. `resolveAbsolutePath("../../examples/tests", cliConfigPath, "")` at line 389
5. `"../"` prefix → `isExplicitRelative = true` → `resolveDotPrefixPath`
6. `source != "runtime"`, `cliConfigPath != ""` → resolves relative to atmos.yaml directory

**Verdict: SAFE.** Identical to previous behavior.

#### Scenario B: `ATMOS_BASE_PATH=.terraform/modules/monorepo` (Spacelift scenario)

1. `processEnvVars` (utils.go:198-203): `BasePath = ".terraform/modules/monorepo"`, `BasePathSource = "runtime"`
2. `resolveAbsolutePath(".terraform/modules/monorepo", cliConfigPath, "runtime")`
3. `.terraform` does NOT match `"./"` — classified as bare path
4. → `tryResolveWithGitRoot(".terraform/modules/monorepo", cliConfigPath)`
5. `os.Stat(gitRoot/.terraform/modules/monorepo)` — doesn't exist at git root
6. Falls back to `os.Stat(CWD/.terraform/modules/monorepo)` — exists → returns CWD path

**Verdict: FIXED.** Previously step 5 returned the non-existent git root path and failed.
Now the `os.Stat` fallback (lines 304-332) correctly finds it at CWD.

#### Scenario C: `ATMOS_BASE_PATH=./.terraform/modules/monorepo` (recommended form)

1. `processEnvVars`: `BasePathSource = "runtime"`
2. `"./"` prefix → `isExplicitRelative = true` → `resolveDotPrefixPath`
3. `source == "runtime"` → `filepath.Abs("./.terraform/modules/monorepo")` = CWD-relative

**Verdict: SAFE.** Direct CWD resolution, no git root lookup needed.

#### Scenario D: Config file `base_path: "../../examples/tests"` (provider test fixture)

Same as Scenario A. Config source, dot-prefixed, resolves relative to atmos.yaml directory.

**Verdict: SAFE.** No behavioral change.

#### Scenario E: Empty `base_path` + run from subdirectory (v1.202.0 feature)

1. `resolveAbsolutePath("", cliConfigPath, "")` → `tryResolveWithGitRoot("", cliConfigPath)`
2. `path == ""` at line 296 → returns git root directly (early return, never hits `os.Stat` code)

**Verdict: SAFE.** Git root discovery preserved exactly as before.

### Risk Assessment

| Path Type                   | Before (v1.209.0)               | After (v1.210.1)                                        | Breaking?  |
|-----------------------------|---------------------------------|---------------------------------------------------------|------------|
| **Empty `""`**              | Git root -> config dir -> CWD   | Same                                                    | No         |
| **Absolute `/abs/path`**    | Pass through                    | Same                                                    | No         |
| **Bare `"foo/bar"`**        | Git root search (no fallback)   | Git root search + `os.Stat` fallback to CWD             | No (safer) |
| **Dot `"./foo"`**           | Resolved relative to config dir | Resolved relative to CWD (runtime source)               | No (fix)   |
| **Relative `"../../path"`** | Resolved relative to config dir | Source-dependent (CWD for runtime, config dir for yaml) | No (fix)   |

### Edge Cases

| Edge Case                                                              | Behavior                                                             | Risk                               |
|------------------------------------------------------------------------|----------------------------------------------------------------------|------------------------------------|
| Both git root and CWD have `.terraform/modules/monorepo`               | Git root wins (line 305 checks first)                                | Low — correct for monorepo layouts |
| `os.Stat` race condition (dir vanishes between stat and use)           | Theoretical only — directories checked are stable (.git, .terraform) | None                               |
| `ATMOS_GIT_ROOT_BASEPATH=false` escape hatch                           | Returns `""` from `getGitRootOrEmpty()` → falls to config path → CWD | None — still works                 |
| `processComponentInStackWithConfig` hardcoded `ProcessTemplates: true` | Not affected — provider already passes `WithProcessTemplates(false)` | None                               |

### Why This Is Safe for Existing Provider Users

1. **Empty `atmos_base_path` (default)** — the most common case — is completely unaffected.
   Empty paths return early at line 296-297 before any new code is reached.

2. **Config-file relative paths** (`../../examples/tests`) — classified as dot-prefixed,
   routed to `resolveDotPrefixPath` with empty source → resolves relative to atmos.yaml dir.
   Unchanged behavior.

3. **The new `os.Stat` fallback** in `tryResolveWithGitRoot` (lines 304-332) is **strictly
   additive** — it only activates when the git-root path doesn't exist, then tries CWD. If
   the git-root path DOES exist, behavior is identical to before. If neither exists, it
   returns the git-root path (line 332) for consistent error messages with pre-fix behavior.

4. **All 7 test packages pass** (15 existing + 4 new tests) with zero regressions.

## Tests

### Existing Tests (Must Pass)

All existing tests in the provider must continue to pass:

| Test File                                              | Tests                                   |
|--------------------------------------------------------|-----------------------------------------|
| `internal/component/component_processor_test.go`       | Component processing, consistency, etc. |
| `internal/describe/describe_stacks_test.go`            | Stack description                       |
| `internal/spacelift/spacelift_stack_processor_test.go` | Spacelift config                        |
| `internal/stack/stack_processor_test.go`               | Stack processing                        |
| `internal/merge/merge_test.go`                         | Deep merge                              |
| `internal/convert/*_test.go`                           | YAML/JSON conversion                    |
| `internal/provider/provider_test.go`                   | Provider schema                         |
| `internal/provider/provider_utils_test.go`             | Provider utilities                      |

### New Tests Added

| Test                                                       | What It Verifies                                                                       |
|------------------------------------------------------------|----------------------------------------------------------------------------------------|
| `TestComponentProcessorWithEmptyBasePath`                  | Empty `atmosBasePath` (default provider behavior) returns correct results              |
| `TestComponentProcessorFromContextWithEmptyBasePath`       | Empty `atmosBasePath` via `ProcessComponentFromContext` returns correct results        |
| `TestComponentProcessorWithRelativeBasePath`               | Relative `base_path` in `atmos.yaml` still works after path resolution changes         |
| `TestComponentProcessorWithProcessingDisabledAndEmptyPath` | Both processing disabled + empty base path (actual provider mode) returns valid config |

### Test Results

All tests pass. Full test suite results:

```text
ok  github.com/cloudposse/terraform-provider-utils/internal/component  4.576s
ok  github.com/cloudposse/terraform-provider-utils/internal/convert    2.203s
ok  github.com/cloudposse/terraform-provider-utils/internal/describe   3.554s
ok  github.com/cloudposse/terraform-provider-utils/internal/merge      1.664s
ok  github.com/cloudposse/terraform-provider-utils/internal/provider   5.392s
ok  github.com/cloudposse/terraform-provider-utils/internal/spacelift  1.327s
ok  github.com/cloudposse/terraform-provider-utils/internal/stack      6.467s
```

New tests (all PASS):

```text
=== RUN   TestComponentProcessorWithEmptyBasePath
--- PASS: TestComponentProcessorWithEmptyBasePath (0.01s)
=== RUN   TestComponentProcessorFromContextWithEmptyBasePath
--- PASS: TestComponentProcessorFromContextWithEmptyBasePath (0.01s)
=== RUN   TestComponentProcessorWithRelativeBasePath
--- PASS: TestComponentProcessorWithRelativeBasePath (0.01s)
=== RUN   TestComponentProcessorWithProcessingDisabledAndEmptyPath
--- PASS: TestComponentProcessorWithProcessingDisabledAndEmptyPath (0.01s)
```

Existing tests (all PASS, no regressions):

```text
--- PASS: TestComponentProcessor (0.13s)
--- PASS: TestComponentProcessorConsistency (0.02s)
--- PASS: TestComponentProcessorProdStack (0.01s)
--- PASS: TestComponentProcessorFromContextProdStack (0.01s)
--- PASS: TestComponentProcessorFromContextNilParams (0.00s)
--- PASS: TestComponentProcessorInfraVpc (0.02s)
--- PASS: TestComponentProcessorWithProcessingDisabled (0.01s)
--- PASS: TestComponentProcessorFromContextWithProcessingDisabled (0.01s)
--- PASS: TestComponentProcessorDisabledMatchesEnabled (0.02s)
--- PASS: TestComponentProcessorConsistencyWithProcessingDisabled (0.02s)
--- PASS: TestComponentProcessorHierarchicalInheritance (0.01s)
--- PASS: TestDescribeStacks (0.14s)
--- PASS: TestDescribeStacksWithFilter1-7 (all pass)
--- PASS: TestMergeBasic and all merge tests (all pass)
--- PASS: TestProvider (schema validation)
--- PASS: TestSpaceliftStackProcessor (all pass)
--- PASS: TestStackProcessor (all pass)
```

## References

- Atmos PR: [cloudposse/atmos#2215](https://github.com/cloudposse/atmos/pull/2215)
- Related issue: [cloudposse/atmos#2183](https://github.com/cloudposse/atmos/issues/2183)

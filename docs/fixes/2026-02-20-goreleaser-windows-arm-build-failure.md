# GoReleaser Fails on `windows/arm` After Go 1.26 Upgrade

**Triggered by:** PR #522 (Atmos v1.207.0 upgrade)

**Affected Versions:** `cloudposse/utils` provider at commit `d668a0c` (post-merge of #522)

**Severity:** High - release pipeline is broken, no new provider versions can be published

## Symptoms

The `release / goreleaser` job in the `Tests` workflow fails after 33 minutes with:

```text
⨯ release failed after 33m24s  error=failed to build for windows_arm_6: exit status 2: go: unsupported GOOS/GOARCH pair windows/arm
```

All 14 binary targets begin building, but the `windows/arm` (32-bit ARM, GOARM=6) target fails
because Go itself rejects the platform pair. The build spends ~33 minutes compiling the other 13
targets before the failure is reported.

CI run: https://github.com/cloudposse/terraform-provider-utils/actions/runs/22229231160/job/64310167511

## Root Cause

PR #522 bumped the Go version from **1.23 to 1.26** (in both `.go-version` and `go.mod`) as
required by the Atmos v1.207.0 dependency.

Go 1.24 (released February 2025) **removed support for the `windows/arm` port** (32-bit ARM on
Windows). From the [Go 1.24 release notes](https://go.dev/doc/go1.24#ports):

> Go 1.24 is the last release that supports building for 386 and arm GOOS targets on Windows.

The provider's release pipeline uses a **shared org-wide goreleaser config** from
`cloudposse/.github/.github/goreleaser.yml`, which builds a full cross-product of:

- **goos:** `freebsd`, `windows`, `linux`, `darwin`
- **goarch:** `amd64`, `386`, `arm`, `arm64`

This produces 16 combinations, but goreleaser silently skips known-invalid pairs (e.g.
`darwin/386`, `darwin/arm`). However, `windows/arm` was a valid pair until Go 1.23, so goreleaser
attempts the build and Go rejects it.

### Build matrix (16 combinations)

| goos    | amd64 | 386  | arm      | arm64 |
|---------|-------|------|----------|-------|
| freebsd | ok    | ok   | ok       | ok    |
| windows | ok    | ok   | **FAIL** | ok    |
| linux   | ok    | ok   | ok       | ok    |
| darwin  | ok    | skip | skip     | ok    |

### Why it worked before

With Go 1.23, `windows/arm` was a supported (though rarely used) build target. The upgrade to
Go 1.26 made it invalid.

### Workflow chain

1. Push to `main` triggers `.github/workflows/test.yml`
2. After `build` and `test` jobs pass, the `release` job calls
   `cloudposse/.github/.github/workflows/shared-go-auto-release.yml@main`
3. The shared workflow checks out `cloudposse/.github` and looks for a local `.goreleaser.yml`:
   ```bash
   if [ -f .goreleaser.yml ]; then
     GORELEASER_CONFIG="./.goreleaser.yml"
   else
     GORELEASER_CONFIG="../.configs/.github/goreleaser.yml"
   fi
   ```
4. No local config existed, so the shared config was used (which includes `windows/arm`)
5. GoReleaser invokes `go build` with `GOOS=windows GOARCH=arm` -> Go exits with
   `unsupported GOOS/GOARCH pair`

## Scope of Impact

This affects **every repo in the `cloudposse` org** that upgrades to Go 1.24+ and uses either:

- The shared org-wide goreleaser config at `cloudposse/.github/.github/goreleaser.yml`
- A local `.goreleaser.yml` that includes `windows` in `goos` and `arm` in `goarch`

### Repos affected

| Repo | goreleaser config | Go version | Status |
|------|-------------------|------------|--------|
| `terraform-provider-utils` | shared (no local config) | 1.26 | **Broken** — failed on release |
| `atmos` | local `.goreleaser.yml` | 1.25.5 | At risk — same `windows/arm` in build matrix |
| Other org repos | shared | varies | Will break when they upgrade to Go 1.24+ |

## Fix

The fix is applied in three places:

### 1. Shared org-wide config (`cloudposse/.github`)

Added an `ignore` rule to `.github/goreleaser.yml` to exclude `windows/arm` from the build matrix:

```yaml
  ignore:
    # Go 1.24+ dropped support for windows/arm (32-bit ARM).
    # https://go.dev/doc/go1.24#ports
    - goos: windows
      goarch: arm
```

This is the primary fix — it prevents all org repos using the shared config from hitting this
failure when they upgrade Go. The `ignore` rule is harmless for repos still on Go < 1.24
(goreleaser simply skips a target that would otherwise build successfully).

Branch: `aknysh/update-goreleaser-1` in `cloudposse/.github`

### 2. Provider: local `.goreleaser.yml` (temporary)

Added a local `.goreleaser.yml` to `terraform-provider-utils` as an immediate workaround,
identical to the shared config but with the `ignore` rule. This file can be **removed** once the
shared config fix in `cloudposse/.github` is merged.

Branch: `aknysh/fix-ci-1`

### 3. Atmos: local `.goreleaser.yml` (needs update)

The `atmos` repo has its own local `.goreleaser.yml` (not using the shared config) and also
includes `windows/arm` in its build matrix. The same `ignore` rule needs to be added there.

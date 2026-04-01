# CLAUDE.md — Maintenance Instructions

This file documents how to maintain `cloudposse/terraform-provider-utils` for automated agents and contributors.

## Releasing a New Version

When a new release is made (by bumping the Atmos dependency in `go.mod`), update the version compatibility
documentation so downstream users can pin correctly.

### Steps

1. **Identify the new Atmos version** from `go.mod`:
   ```
   grep 'cloudposse/atmos' go.mod
   ```

2. **Update `docs/versions.md`** — add a new row at the top of the table:
   ```markdown
   | [vX.Y.Z](https://github.com/cloudposse/atmos/releases/tag/vX.Y.Z) | [vA.B.C](https://github.com/cloudposse/terraform-provider-utils/releases/tag/vA.B.C) |
   ```
   Replace `vX.Y.Z` with the Atmos version and `vA.B.C` with the new provider release tag.

   - Only add one row per unique Atmos version.
   - If the new provider release uses the same Atmos version as the previous one, **update** the existing row
     to point to the newer provider tag instead of adding a duplicate.
   - If there are parallel `1.x` and `2.x` releases for the same Atmos version, prefer the `2.x` release as
     the recommended version.

3. **Update `README.md`** — the version table in `README.md` is a copy of the table in `docs/versions.md`
   (it is inlined when `atmos docs generate` is run). Replicate the same change in `README.md` under the
   **"Version Compatibility"** section so the README stays in sync without requiring a full docs regeneration.

4. **Commit with a semantic commit message**, e.g.:
   ```
   docs: update version compatibility matrix for atmos vX.Y.Z / provider vA.B.C
   ```

### Why Version Pinning Matters

The `cloudposse/utils` provider embeds the Atmos Go library to perform YAML/JSON deep merging.
If the embedded library version differs from the Atmos CLI version used to orchestrate stacks, the two
tools may resolve deep-merge conflicts differently, producing non-deterministic infrastructure state.
Always keep both at matching versions.

## Regenerating README.md

The canonical source for the README is `README.yaml`. After editing `README.yaml` or any file in `docs/`,
regenerate `README.md` with:

```sh
atmos docs generate
```

If you cannot run atmos locally, manually mirror changes from `docs/versions.md` into the
**"Version Compatibility"** section of `README.md`.

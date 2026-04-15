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

2. **Update `docs/versions.md`** — add or update a row in the table.

   The table has three columns: **Atmos Version**, **2.x provider release**, **1.x provider release**.

   - If the new release introduces a **new Atmos version**, add a new row at the top of the table:
     ```markdown
     | [vX.Y.Z](https://github.com/cloudposse/atmos/releases/tag/vX.Y.Z) | [v2.A.B](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v2.A.B) | — |
     ```
     (Use `—` in the column that does not apply to this release.)
   - If the new release uses the **same Atmos version** as the previous one, **update** the existing row
     to point to the newer provider tag in the appropriate column instead of adding a duplicate row.
   - The `2.x` column tracks the actively developed release line; the `1.x` column tracks parallel
     backport releases for users who cannot yet migrate to `2.x`.

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

## Updating the `terraform-yaml-stack-config` Compatibility Table

`docs/versions.md` (and its mirror in `README.md`) also contains a second table documenting which releases of
[`cloudposse/terraform-yaml-stack-config`](https://github.com/cloudposse/terraform-yaml-stack-config) require
which `cloudposse/utils` version. This matters because the `modules/remote-state` child module inside that repo
declares explicit provider version constraints that must be satisfiable together with the root module's constraint.

When a new `terraform-yaml-stack-config` release changes its `cloudposse/utils` version constraint (check
`versions.tf` and `modules/remote-state/versions.tf` in that repo), **add or update a row** in the
`terraform-yaml-stack-config` ↔ `cloudposse/utils` table in both `docs/versions.md` and `README.md`.

If the new release introduces a new constraint range, add a new row. If it stays within an existing range,
update the upper bound of the existing range row (e.g. change `v1.5.0–v1.8.0` to `v1.5.0–v1.9.0`).

## Regenerating README.md

The canonical source for the README is `README.yaml`. After editing `README.yaml` or any file in `docs/`,
regenerate `README.md` with:

```sh
atmos docs generate
```

If you cannot run atmos locally, manually mirror changes from `docs/versions.md` into the
**"Version Compatibility"** section of `README.md`.

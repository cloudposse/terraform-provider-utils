## Version Compatibility: Atmos ↔ `cloudposse/utils` Provider

The `cloudposse/utils` Terraform provider embeds the [Atmos](https://atmos.tools) Go library to perform
[deep-merge](https://atmos.tools/core-concepts/stacks/inheritance/deep-merging/) operations on YAML and JSON.
Because deep-merge semantics can change between Atmos releases, **the Atmos library version compiled into the
provider must match the Atmos CLI version you are running.** A mismatch can cause the provider and Atmos to
produce different merge results for the same inputs, breaking deterministic infrastructure configuration.

Use the table below to pin both tools to compatible versions.

> **Note:** The `2.x` line is the actively developed release line. The `1.x` line receives parallel backport
> releases for users who cannot yet migrate. Where multiple provider releases share an Atmos version, only the
> latest release in each line is shown.

| Atmos Version | `cloudposse/utils` 2.x | `cloudposse/utils` 1.x |
|---|---|---|
| [v1.211.0](https://github.com/cloudposse/atmos/releases/tag/v1.211.0) | [v2.5.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v2.5.0) | [v1.35.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v1.35.0) |
| [v1.210.1](https://github.com/cloudposse/atmos/releases/tag/v1.210.1) | [v2.4.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v2.4.0) | [v1.34.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v1.34.0) |
| [v1.209.0](https://github.com/cloudposse/atmos/releases/tag/v1.209.0) | [v2.1.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v2.1.0) | [v1.33.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v1.33.0) |
| [v1.207.0](https://github.com/cloudposse/atmos/releases/tag/v1.207.0) | [v2.0.2](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v2.0.2) | — |
| [v1.189.0](https://github.com/cloudposse/atmos/releases/tag/v1.189.0) | — | [v1.31.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v1.31.0) |
| [v1.177.0](https://github.com/cloudposse/atmos/releases/tag/v1.177.0) | — | [v1.30.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v1.30.0) |
| [v1.165.3](https://github.com/cloudposse/atmos/releases/tag/v1.165.3) | — | [v1.29.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v1.29.0) |
| [v1.130.0](https://github.com/cloudposse/atmos/releases/tag/v1.130.0) | — | [v1.28.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v1.28.0) |
| [v1.122.0](https://github.com/cloudposse/atmos/releases/tag/v1.122.0) | — | [v1.27.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v1.27.0) |
| [v1.88.0](https://github.com/cloudposse/atmos/releases/tag/v1.88.0) | — | [v1.26.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v1.26.0) |
| [v1.87.0](https://github.com/cloudposse/atmos/releases/tag/v1.87.0) | — | [v1.25.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v1.25.0) |
| [v1.84.0](https://github.com/cloudposse/atmos/releases/tag/v1.84.0) | — | [v1.24.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/v1.24.0) |
| [v1.79.0](https://github.com/cloudposse/atmos/releases/tag/v1.79.0) | — | [v1.23.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.23.0) |
| [v1.70.0](https://github.com/cloudposse/atmos/releases/tag/v1.70.0) | — | [v1.22.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.22.0) |
| [v1.68.0](https://github.com/cloudposse/atmos/releases/tag/v1.68.0) | — | [v1.21.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.21.0) |
| [v1.66.0](https://github.com/cloudposse/atmos/releases/tag/v1.66.0) | — | [v1.19.2](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.19.2) |
| [v1.64.1](https://github.com/cloudposse/atmos/releases/tag/v1.64.1) | — | [v1.19.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.19.0) |
| [v1.63.0](https://github.com/cloudposse/atmos/releases/tag/v1.63.0) | — | [v1.17.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.17.0) |
| [v1.54.0](https://github.com/cloudposse/atmos/releases/tag/v1.54.0) | — | [v1.16.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.16.0) |
| [v1.50.0](https://github.com/cloudposse/atmos/releases/tag/v1.50.0) | — | [v1.14.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.14.0) |
| [v1.49.0](https://github.com/cloudposse/atmos/releases/tag/v1.49.0) | — | [v1.13.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.13.0) |
| [v1.45.0](https://github.com/cloudposse/atmos/releases/tag/v1.45.0) | — | [v1.12.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.12.0) |
| [v1.44.0](https://github.com/cloudposse/atmos/releases/tag/v1.44.0) | — | [v1.11.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.11.0) |
| [v1.43.0](https://github.com/cloudposse/atmos/releases/tag/v1.43.0) | — | [v1.10.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.10.0) |
| [v1.41.0](https://github.com/cloudposse/atmos/releases/tag/v1.41.0) | — | [v1.9.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.9.0) |
| [v1.34.2](https://github.com/cloudposse/atmos/releases/tag/v1.34.2) | — | [v1.8.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.8.0) |
| [v1.26.1](https://github.com/cloudposse/atmos/releases/tag/v1.26.1) | — | [v1.7.1](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.7.1) |
| [v1.26.0](https://github.com/cloudposse/atmos/releases/tag/v1.26.0) | — | [v1.7.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.7.0) |
| [v1.15.0](https://github.com/cloudposse/atmos/releases/tag/v1.15.0) | — | [v1.6.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.6.0) |
| [v1.10.1](https://github.com/cloudposse/atmos/releases/tag/v1.10.1) | — | [v1.5.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.5.0) |
| [v1.10.0](https://github.com/cloudposse/atmos/releases/tag/v1.10.0) | — | [v1.4.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.4.0) |
| [v1.9.1](https://github.com/cloudposse/atmos/releases/tag/v1.9.1) | — | [v1.4.1](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.4.1) |
| [v1.8.1](https://github.com/cloudposse/atmos/releases/tag/v1.8.1) | — | [v1.2.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.2.0) |
| [v1.8.0](https://github.com/cloudposse/atmos/releases/tag/v1.8.0) | — | [v1.1.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.1.0) |
| [v1.4.26](https://github.com/cloudposse/atmos/releases/tag/v1.4.26) | — | [v1.0.0](https://github.com/cloudposse/terraform-provider-utils/releases/tag/1.0.0) |
| [v1.4.25](https://github.com/cloudposse/atmos/releases/tag/v1.4.25) | — | [v0.17.28](https://github.com/cloudposse/terraform-provider-utils/releases/tag/0.17.28) |
| [v1.4.24](https://github.com/cloudposse/atmos/releases/tag/v1.4.24) | — | [v0.17.27](https://github.com/cloudposse/terraform-provider-utils/releases/tag/0.17.27) |
| [v1.4.22](https://github.com/cloudposse/atmos/releases/tag/v1.4.22) | — | [v0.17.26](https://github.com/cloudposse/terraform-provider-utils/releases/tag/0.17.26) |
| [v1.4.20](https://github.com/cloudposse/atmos/releases/tag/v1.4.20) | — | [v0.17.25](https://github.com/cloudposse/terraform-provider-utils/releases/tag/0.17.25) |
| [v1.4.15](https://github.com/cloudposse/atmos/releases/tag/v1.4.15) | — | [v0.17.24](https://github.com/cloudposse/terraform-provider-utils/releases/tag/0.17.24) |
| [v1.4.11](https://github.com/cloudposse/atmos/releases/tag/v1.4.11) | — | [v0.17.23](https://github.com/cloudposse/terraform-provider-utils/releases/tag/0.17.23) |
| [v1.4.10](https://github.com/cloudposse/atmos/releases/tag/v1.4.10) | — | [v0.17.22](https://github.com/cloudposse/terraform-provider-utils/releases/tag/0.17.22) |
| [v1.4.8](https://github.com/cloudposse/atmos/releases/tag/v1.4.8) | — | [v0.17.21](https://github.com/cloudposse/terraform-provider-utils/releases/tag/0.17.21) |
| [v1.4.5](https://github.com/cloudposse/atmos/releases/tag/v1.4.5) | — | [v0.17.20](https://github.com/cloudposse/terraform-provider-utils/releases/tag/0.17.20) |
| [v1.4.2](https://github.com/cloudposse/atmos/releases/tag/v1.4.2) | — | [v0.17.19](https://github.com/cloudposse/terraform-provider-utils/releases/tag/0.17.19) |
| [v1.4.1](https://github.com/cloudposse/atmos/releases/tag/v1.4.1) | — | [v0.17.18](https://github.com/cloudposse/terraform-provider-utils/releases/tag/0.17.18) |

## Version Compatibility: `terraform-yaml-stack-config` ↔ `cloudposse/utils` Provider

The [`cloudposse/terraform-yaml-stack-config`](https://github.com/cloudposse/terraform-yaml-stack-config) module
uses the `cloudposse/utils` provider internally. Its `modules/remote-state` child module declares explicit provider
version constraints — if you source that child module in a root module, both the child's constraint and any
constraint declared by the root module must be satisfiable simultaneously.

Use the table below to pick a `terraform-yaml-stack-config` release that is compatible with your chosen
`cloudposse/utils` provider version.

| [`terraform-yaml-stack-config`](https://github.com/cloudposse/terraform-yaml-stack-config/releases) | Required `cloudposse/utils` (root module) | Required `cloudposse/utils` (`modules/remote-state`) |
|---|---|---|
| [v2.0.0](https://github.com/cloudposse/terraform-yaml-stack-config/releases/tag/v2.0.0) | `>= 2.0.0, < 3.0.0` | `>= 2.0.0, < 3.0.0` |
| [v1.5.0](https://github.com/cloudposse/terraform-yaml-stack-config/releases/tag/1.5.0)–[v1.8.0](https://github.com/cloudposse/terraform-yaml-stack-config/releases/tag/v1.8.0) | `>= 1.7.1` | `>= 1.7.1, < 2.0.0` |
| [v1.1.1](https://github.com/cloudposse/terraform-yaml-stack-config/releases/tag/1.1.1)–[v1.4.3](https://github.com/cloudposse/terraform-yaml-stack-config/releases/tag/1.4.3) | `>= 1.5.0` | `= 1.5.0` |
| [v1.0.0](https://github.com/cloudposse/terraform-yaml-stack-config/releases/tag/1.0.0) | `>= 1.2.0` | `>= 1.2.0` |

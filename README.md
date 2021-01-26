# Terraform Provider Cloud Posse Utils

TODO: Update this

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

### Testing Locally

You can test the provider locally by using the [provider_installation](https://www.terraform.io/docs/cli/config/config-file.html#provider-installation) functionality.

In our case, you can edit your `~/.terraformrc` file with the following:

```hcl
provider_installation {
  dev_overrides  {
    "cloudposse/utils" = "/path/to/your/code/github.com/cloudposse/terraform-provider-utils/"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

With that in place, you can build the provider (see above) and add a provider block:

```hcl
required_providers {
    utils = {
      source = "cloudposse/utils"
    }
  }
```

Then run `terraform plan` and `terraform apply` as normal. NOTE: `terraform init` will give an error similar to the one below, but per the documentation, that is to be expected and can be ignored.

```sh
$ terraform init
Initializing the backend...

Initializing provider plugins...
- Finding latest version of cloudposse/utils...

Warning: Provider development overrides are in effect

The following provider development overrides are set in the CLI configuration:
 - cloudposse/utils in /Users/matt/code/src/github.com/cloudposse/terraform-provider-utils

The behavior may therefore not match any released version of the provider and
applying changes may cause the state to become incompatible with published
releases.


Error: Failed to query available provider packages

Could not retrieve the list of available versions for provider
cloudposse/utils: provider registry registry.terraform.io does not have a
provider named registry.terraform.io/cloudposse/utils
```

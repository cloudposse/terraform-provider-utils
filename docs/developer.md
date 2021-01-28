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

For testing this provider, you can edit your `~/.terraformrc` file with the following:

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

Then run `terraform init`, `terraform plan` and `terraform apply` as normal.

```sh
$ terraform init
Initializing the backend...

Initializing provider plugins...
- Finding latest version of cloudposse/utils...

Warning: Provider development overrides are in effect

The following provider development overrides are set in the CLI configuration:
 - cloudposse/utils in /path/to/your/code/github.com/cloudposse/terraform-provider-utils

The behavior may therefore not match any released version of the provider and
applying changes may cause the state to become incompatible with published
releases.
```

```sh
terraform apply

Warning: Provider development overrides are in effect

The following provider development overrides are set in the CLI configuration:
 - cloudposse/utils in /Users/matt/code/src/github.com/cloudposse/terraform-provider-utils

The behavior may therefore not match any released version of the provider and
applying changes may cause the state to become incompatible with published
releases.


An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:

Terraform will perform the following actions:

Plan: 0 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + deep_merge_output = <<-EOT
        Statement:
        - Action:
          - s3:*
          Effect: Allow
          Resource:
          - '*'
          Sid: FullAccess
        - Action:
          - s3:*
          Complex:
            ExtraComplex:
              ExtraExtraComplex:
                Foo: bazzz
                SomeArray:
                - one
                - two
                - three
          Effect: Deny
          Resource:
          - arn:aws:s3:::customer
          - arn:aws:s3:::customer/*
          - foo
          Sid: DenyCustomerBucket
        Version: "2012-10-17"
    EOT
```

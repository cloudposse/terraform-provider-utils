terraform {
  required_providers {
    # Install the provider on local computer by running `make install` from the root of the repo
    utils = {
      source  = "cloudposse/utils"
      version = "9999.99.99"
    }
  }
}

data "utils_stack_config_yaml" "example" {
  input = [
    "${path.module}/stacks/uw2-dev.yaml",
    "${path.module}/stacks/uw2-prod.yaml",
    "${path.module}/stacks/uw2-staging.yaml",
    "${path.module}/stacks/uw2-uat.yaml"
  ]
}

output "output" {
  value = data.utils_stack_config_yaml.example.output
}

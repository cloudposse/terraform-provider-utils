terraform {
  required_providers {
    utils = {
      source = "cloudposse/utils"
      # For local development,
      # install the provider on local computer by running `make install` from the root of the repo,
      # and uncomment the version below
      version = "9999.99.99"
    }
  }
}

locals {
  base_path = "../../config/stacks"

  stack_config_files = [
    "${local.base_path}/tenant1/ue2/dev.yaml",
    "${local.base_path}/tenant1/ue2/prod.yaml",
    "${local.base_path}/tenant1/ue2/staging.yaml",
    "${local.base_path}/tenant2/ue2/dev.yaml",
    "${local.base_path}/tenant2/ue2/prod.yaml",
    "${local.base_path}/tenant2/ue2/staging.yaml",
  ]
}

data "utils_spacelift_stack_config" "example" {
  base_path                  = local.base_path
  input                      = local.stack_config_files
  process_stack_deps         = false
  process_component_deps     = true
  process_imports            = true
  stack_config_path_template = "stacks/%s.yaml"
}

locals {
  result = yamldecode(data.utils_spacelift_stack_config.example.output)
}

output "output" {
  value = local.result
}

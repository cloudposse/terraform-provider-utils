locals {
  base_path = "../../complete/stacks"

  stack_config_files = [
    "${local.base_path}/tenant1/ue2/dev.yaml",
    "${local.base_path}/tenant1/ue2/prod.yaml",
    "${local.base_path}/tenant1/ue2/staging.yaml",
    "${local.base_path}/tenant2/ue2/dev.yaml",
    "${local.base_path}/tenant2/ue2/prod.yaml",
    "${local.base_path}/tenant2/ue2/staging.yaml",
  ]

  result = [for i in data.utils_stack_config_yaml.example.output : yamldecode(i)]
}

data "utils_stack_config_yaml" "example" {
  base_path              = local.base_path
  input                  = local.stack_config_files
  process_component_deps = true
  process_stack_deps     = false
}

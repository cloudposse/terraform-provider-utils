locals {
  base_path = "../../tests/stacks"

  stack_config_files = [
    "${local.base_path}/tenant1/ue2/dev.yaml",
    "${local.base_path}/tenant1/ue2/prod.yaml",
    "${local.base_path}/tenant1/ue2/staging.yaml",
    "${local.base_path}/tenant2/ue2/dev.yaml",
    "${local.base_path}/tenant2/ue2/prod.yaml",
    "${local.base_path}/tenant2/ue2/staging.yaml",
  ]

  result = yamldecode(data.utils_spacelift_stack_config.example.output)
}

data "utils_spacelift_stack_config" "example" {
  base_path                  = local.base_path
  input                      = local.stack_config_files
  process_component_deps     = true
  process_stack_deps         = false
  process_imports            = false
  stack_config_path_template = "stacks/%s.yaml"
}

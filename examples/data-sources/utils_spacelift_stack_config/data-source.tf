locals {
  result = yamldecode(data.utils_spacelift_stack_config.example.output)
}

data "utils_spacelift_stack_config" "example" {
  process_component_deps     = true
  stack_config_path_template = "stacks/%s.yaml"
}

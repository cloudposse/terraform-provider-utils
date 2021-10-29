locals {
  component = "test/test-component-override"
  stack     = "tenant1-ue2-dev"

  result = yamldecode(data.utils_component_config.example.output)
}

data "utils_component_config" "example" {
  component = local.component
  stack     = local.stack
}

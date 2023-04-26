locals {
  stack       = "tenant1-ue2-dev"
  namespace   = ""
  tenant      = "tenant1"
  environment = "ue2"
  stage       = "dev"

  env = {
    ENVIRONMENT           = local.environment
    STAGE                 = local.stage
    ATMOS_CLI_CONFIG_PATH = "."
  }

  result1 = yamldecode(data.utils_describe_stacks.example1.output)
  result2 = yamldecode(data.utils_describe_stacks.example2.output)
  result3 = yamldecode(data.utils_describe_stacks.example3.output)
  result4 = yamldecode(data.utils_describe_stacks.example4.output)
  result5 = yamldecode(data.utils_describe_stacks.example5.output)
  result6 = yamldecode(data.utils_describe_stacks.example6.output)
  result7 = yamldecode(data.utils_describe_stacks.example7.output)
}

data "utils_describe_stacks" "example1" {
}

data "utils_describe_stacks" "example2" {
  stack = local.stack
}

data "utils_describe_stacks" "example3" {
  namespace   = local.namespace
  tenant      = local.tenant
  environment = local.environment
  stage       = local.stage
  env         = local.env
}

data "utils_describe_stacks" "example4" {
  atmos_cli_config_path = "."
  atmos_base_path       = "../../complete"
  component_types       = ["terraform"]
  sections              = ["none"]
}

data "utils_describe_stacks" "example5" {
  component_types = ["terraform"]
  components      = ["top-level-component1", "test/test-component-override-3"]
  sections        = ["none"]
}

data "utils_describe_stacks" "example6" {
  tenant          = local.tenant
  environment     = local.environment
  stage           = local.stage
  component_types = ["terraform"]
  components      = ["test/test-component-override-3"]
  sections        = ["vars", "metadata", "env"]
}

data "utils_describe_stacks" "example7" {
  component_types = ["terraform"]
  components      = ["top-level-component1"]
  sections        = ["vars"]
}

locals {
  component   = "test/test-component-override"
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

  result1 = yamldecode(data.utils_component_config.example1.output)
  result2 = yamldecode(data.utils_component_config.example2.output)
}

data "utils_component_config" "example1" {
  component     = local.component
  stack         = local.stack
  ignore_errors = false
  env           = local.env
}

data "utils_component_config" "example2" {
  component     = local.component
  namespace     = local.namespace
  tenant        = local.tenant
  environment   = local.environment
  stage         = local.stage
  ignore_errors = false
  env           = local.env
}

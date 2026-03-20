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
  result3 = yamldecode(data.utils_component_config.example3.output)
  result4 = yamldecode(data.utils_component_config.example4.output)
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

data "utils_component_config" "example3" {
  component             = local.component
  stack                 = local.stack
  ignore_errors         = false
  atmos_cli_config_path = "."
  atmos_base_path       = "../../tests"
}

# Disable Go template processing (enabled by default).
# YAML function processing (e.g., !terraform.output) is disabled by default.
data "utils_component_config" "example4" {
  component              = local.component
  stack                  = local.stack
  ignore_errors          = false
  env                    = local.env
  process_templates      = false
  process_yaml_functions = false
}

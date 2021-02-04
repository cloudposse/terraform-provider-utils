terraform {
  required_providers {
    utils = {
      source  = "cloudposse/utils"
      version = "9999.99.99"
    }
  }
}

locals {
  yaml_data_1 = file("${path.module}/data1.yaml")
  yaml_data_2 = file("${path.module}/data2.yaml")
}

data "utils_stack_config_yaml" "example" {
  inputs = [
    local.yaml_data_1,
    local.yaml_data_2
  ]
}

output "deep_merge_output" {
  value = data.utils_stack_config_yaml.example.output
}

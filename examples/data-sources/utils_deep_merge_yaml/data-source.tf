terraform {
  required_providers {
    utils = {
      source = "cloudposse/utils"
      # Install the provider on local computer by running `make install` from the root of the repo
      # version = "9999.99.99"
    }
  }
}

locals {
  yaml_data_1 = file("${path.module}/data1.yaml")
  yaml_data_2 = file("${path.module}/data2.yaml")
}

data "utils_deep_merge_yaml" "example" {
  input = [
    local.yaml_data_1,
    local.yaml_data_2
  ]
}

output "deep_merge_output" {
  value = yamldecode(data.utils_deep_merge_yaml.example.output)
}

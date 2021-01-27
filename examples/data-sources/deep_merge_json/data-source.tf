terraform {
  required_providers {
    utils = {
      source = "cloudposse/utils"
    }
  }
}

locals {
  json_data_1 = file("${path.module}/json1.json")
  json_data_2 = file("${path.module}/json2.json")
}

data "utils_deep_merge_json" "example" {
  inputs = [
    local.json_data_1,
    local.json_data_2
  ]
}

output "deep_merge_output" {
  value = data.utils_deep_merge_json.example.output
}

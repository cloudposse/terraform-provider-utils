locals {
  json_data_1 = file("${path.module}/json1.json")
  json_data_2 = file("${path.module}/json2.json")
}

data "utils_deep_merge_json" "example" {
  input = [
    local.json_data_1,
    local.json_data_2
  ]
}

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

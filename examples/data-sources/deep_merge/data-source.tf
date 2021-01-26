terraform {
  required_providers {
    utils = {
      source = "cloudposse/utils"
    }
  }
}

data "utils_deep_merge" "example" {
  inputs = [
    { foo = "bar" },
    { baz = "bat" },
    { foo = "ood" },
    { baz = "zzed" },
    { my = "dinner" }
  ]
}

output "deep_merge_output" {
  value = data.utils_deep_merge.example.output
}

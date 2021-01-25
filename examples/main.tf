terraform {
  required_providers {
    cloudposse = {
      source = "cloudposse/utils"
    }
  }
}


data "cloudposse_deep_merge" "dm1" {
  inputs = {
    foo = "bar",
    baz = "bat"
  }
}

output "test" {
  value = data.cloudposse_deep_merge.dm1.output
}

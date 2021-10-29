output "output" {
  value = local.result
}

output "tenant1_ue2_dev_echo_server_vars" {
  value = local.result[0]["components"]["helmfile"]["echo-server"]["vars"]
}

output "tenant1_ue2_dev_infra_vpc_vars" {
  value = local.result[0]["components"]["terraform"]["infra/vpc"]["vars"]
}

output "tenant1_ue2_dev_test_component_vars" {
  value = local.result[0]["components"]["terraform"]["test/test-component"]["vars"]
}

output "tenant1_ue2_dev_test_component_override_vars" {
  value = local.result[0]["components"]["terraform"]["test/test-component-override"]["vars"]
}

output "tenant1_ue2_dev_test_component_override_backend" {
  value = local.result[0]["components"]["terraform"]["test/test-component-override"]["backend"]
}

output "tenant1_ue2_dev_test_component_override_deps" {
  value = local.result[0]["components"]["terraform"]["test/test-component-override"]["deps"]
}

output "tenant1_ue2_dev_test_component_override_settings" {
  value = local.result[0]["components"]["terraform"]["test/test-component-override"]["settings"]
}

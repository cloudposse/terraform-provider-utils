output "output1" {
  value       = local.result1
  description = "Component config from provided stack"
}

output "output2" {
  value       = local.result2
  description = "Component config from provided context"
}

output "output3" {
  value       = local.result3
  description = "Component config from provided stack using `atmos_cli_config_path` and `atmos_base_path` variables"
}

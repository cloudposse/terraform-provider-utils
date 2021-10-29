output "output1" {
  value       = local.result1
  description = "Component config from provided stack"
}

output "output2" {
  value       = local.result2
  description = "Component config from provided context"
}

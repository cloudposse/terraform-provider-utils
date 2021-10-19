output "deep_merge_output" {
  value = yamldecode(data.utils_deep_merge_yaml.example.output)
}

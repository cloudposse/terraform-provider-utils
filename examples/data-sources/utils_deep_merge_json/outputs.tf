output "deep_merge_output" {
  value = jsondecode(data.utils_deep_merge_json.example.output)
}

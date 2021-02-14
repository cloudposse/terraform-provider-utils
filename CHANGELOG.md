## 0.3.0

- Workaround for a bug in `mergo.Merge()` when merging slice of maps in a `for` loop, 
  it modifies the source of the previous loop iteration if it's a complex map and it get a pointer to it, 
  not only the destination of the current loop iteration)

- Added `settings` section to YAML stack configurations for Terraform and helmfile components

- Added `env` section to YAML stack configurations for Terraform and helmfile components

BACKWARDS INCOMPATIBILITIES / NOTES:


## 0.2.0

- Added `data_source_stack_config_yaml` data source to process YAML stack configurations for Terraform and helmfile components

BACKWARDS INCOMPATIBILITIES / NOTES:


## 0.1.0

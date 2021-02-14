## 0.3.0

- Workaround for a deep-merge bug in `mergo.Merge()`. When deep-merging slice of maps in a `for` loop, 
  `mergo` modifies the source of the previous loop iteration if it's a complex map and `mergo` gets a pointer to it, 
  not only the destination of the current loop iteration.

- Added `settings` sections to `data_source_stack_config_yaml` data source to provide settings for Terraform and helmfile components

- Added `env` sections to `data_source_stack_config_yaml` data source to provide ENV vars for Terraform and helmfile components

BACKWARDS INCOMPATIBILITIES / NOTES:


## 0.2.0

- Added `data_source_stack_config_yaml` data source to process YAML stack configurations for Terraform and helmfile components

BACKWARDS INCOMPATIBILITIES / NOTES:


## 0.1.0

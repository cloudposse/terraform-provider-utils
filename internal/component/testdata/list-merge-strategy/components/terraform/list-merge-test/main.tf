# Minimal component used only to give the stack manifests a real Terraform
# component directory to point at. The provider/Atmos only processes the stack
# configuration here; these variables are never applied.
variable "my_list" {
  type    = any
  default = []
}

variable "my_map_list" {
  type    = any
  default = []
}

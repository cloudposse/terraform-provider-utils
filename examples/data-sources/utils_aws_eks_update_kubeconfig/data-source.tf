/*
These examples show how to use the 'utils_aws_eks_update_kubeconfig' data source.

The data source can executes 'aws eks update-kubeconfig' commands in four different ways:

1. If all the required parameters (cluster name and AWS profile/role) are provided,
   then it executes the command without requiring the CLI config ('atmos.yaml') and component/stack/context.
   'atmos.yaml' is not required/needed in this case.
   See 'example1'.

2. If 'component' and 'stack' are provided,
   then it executes the command using the atmos CLI config (see atmos.yaml) and the context by searching for the following settings:
     - 'components.helmfile.cluster_name_pattern' in 'atmos.yaml' CLI config (and calculates the '--name' parameter using the pattern)
     - 'components.helmfile.helm_aws_profile_pattern' in 'atmos.yaml' CLI config (and calculates the '--profile' parameter using the pattern)
     - 'components.helmfile.kubeconfig_path' in 'atmos.yaml' CLI config
     - the variables for the component in the provided stack
     - 'region' from the variables for the component in the stack
  See 'example2'.

3. If the context ('tenant', 'environment', 'stage') and 'component' are provided,
   then it builds the stack name by using the 'stacks.name_pattern' CLI config from 'atmos.yaml', then performs the same steps as example #2.
   See 'example3'.

4. Combination of the above. Provide a component and a stack (or context), and override other parameters (e.g. 'kubeconfig', 'region').
   See 'example4'.

If 'kubeconfig' (the filename to write the kubeconfig to) is not provided, then it's calculated by joining
the base path from 'components.helmfile.kubeconfig_path' CLI config from 'atmos.yaml' and the stack name.

Supported inputs of the 'utils_aws_eks_update_kubeconfig' data source:
  - component
  - stack
  - tenant
  - environment
  - stage
  - cluster_name
  - kubeconfig
  - profile
  - role_arn
  - alias
  - region

Docs: https://docs.aws.amazon.com/cli/latest/reference/eks/update-kubeconfig.html
*/

locals {
  component   = "test/test-component-override"
  stack       = "tenant1-ue2-dev"
  tenant      = "tenant1"
  environment = "ue2"
  stage       = "dev"

  profile      = "eg-gbl-dev-admin"
  cluster_name = "eg-ue2-dev-eks-cluster"
  kubeconfig   = "./kubeconfig"
  region       = "us-east-2"

  result1 = data.utils_aws_eks_update_kubeconfig.example1.output
  result2 = data.utils_aws_eks_update_kubeconfig.example2.output
  result3 = data.utils_aws_eks_update_kubeconfig.example3.output
  result4 = data.utils_aws_eks_update_kubeconfig.example4.output
}

data "utils_aws_eks_update_kubeconfig" "example1" {
  profile      = local.profile
  cluster_name = local.cluster_name
  kubeconfig   = local.kubeconfig
}

data "utils_aws_eks_update_kubeconfig" "example2" {
  component  = local.component
  stack      = local.stack
  kubeconfig = local.kubeconfig
}

data "utils_aws_eks_update_kubeconfig" "example3" {
  component   = local.component
  tenant      = local.tenant
  environment = local.environment
  stage       = local.stage
}

data "utils_aws_eks_update_kubeconfig" "example4" {
  component   = local.component
  tenant      = local.tenant
  environment = local.environment
  stage       = local.stage
  kubeconfig  = local.kubeconfig
  region      = local.region
}

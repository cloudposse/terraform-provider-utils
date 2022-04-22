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
  component = local.component
  stack     = local.stack
}

data "utils_aws_eks_update_kubeconfig" "example2" {
  component   = local.component
  tenant      = local.tenant
  environment = local.environment
  stage       = local.stage
}

data "utils_aws_eks_update_kubeconfig" "example3" {
  profile      = local.profile
  cluster_name = local.cluster_name
  kubeconfig   = local.kubeconfig
}

data "utils_aws_eks_update_kubeconfig" "example4" {
  component   = local.component
  tenant      = local.tenant
  environment = local.environment
  stage       = local.stage
  region      = local.region
}

package spacelift

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestSpaceliftStackProcessorWithTerraform(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "./",
		Upgrade:      true,
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	var output interface{}
	terraform.OutputStruct(t, terraformOptions, "output", &output)
	spaceliftStacks := output.(map[string]interface{})

	assert.Equal(t, 30, len(spaceliftStacks))

	tenant1Ue2DevInfraVpcStack := spaceliftStacks["tenant1-ue2-dev-infra-vpc"].(map[string]interface{})
	tenant1Ue2DevInfraVpcStackInfrastructureStackName := tenant1Ue2DevInfraVpcStack["stack"].(string)
	tenant1Ue2DevInfraVpcStackBackend := tenant1Ue2DevInfraVpcStack["backend"].(map[string]interface{})
	tenant1Ue2DevInfraVpcStackBackendWorkspaceKeyPrefix := tenant1Ue2DevInfraVpcStackBackend["workspace_key_prefix"].(string)
	assert.Equal(t, "tenant1-ue2-dev", tenant1Ue2DevInfraVpcStackInfrastructureStackName)
	assert.Equal(t, "infra-vpc", tenant1Ue2DevInfraVpcStackBackendWorkspaceKeyPrefix)

	tenant1Ue2DevTestTestComponentOverrideComponent := spaceliftStacks["tenant1-ue2-dev-test-test-component-override"].(map[string]interface{})
	tenant1Ue2DevTestTestComponentOverrideComponentInfrastructureStackName := tenant1Ue2DevTestTestComponentOverrideComponent["stack"].(string)
	tenant1Ue2DevTestTestComponentOverrideComponentBackend := tenant1Ue2DevTestTestComponentOverrideComponent["backend"].(map[string]interface{})
	tenant1Ue2DevTestTestComponentOverrideComponentBaseComponent := tenant1Ue2DevTestTestComponentOverrideComponent["base_component"].(string)
	tenant1Ue2DevTestTestComponentOverrideComponentBackendWorkspaceKeyPrefix := tenant1Ue2DevTestTestComponentOverrideComponentBackend["workspace_key_prefix"].(string)
	tenant1Ue2DevTestTestComponentOverrideComponentDeps := tenant1Ue2DevTestTestComponentOverrideComponent["deps"].([]interface{})
	tenant1Ue2DevTestTestComponentOverrideComponentLabels := tenant1Ue2DevTestTestComponentOverrideComponent["labels"].([]interface{})
	tenant1Ue2DevTestTestComponentOverrideTerraformWorkspace := tenant1Ue2DevTestTestComponentOverrideComponent["workspace"]
	assert.Equal(t, "tenant1-ue2-dev", tenant1Ue2DevTestTestComponentOverrideComponentInfrastructureStackName)
	assert.Equal(t, "test-test-component", tenant1Ue2DevTestTestComponentOverrideComponentBackendWorkspaceKeyPrefix)
	assert.Equal(t, "test/test-component", tenant1Ue2DevTestTestComponentOverrideComponentBaseComponent)
	assert.Equal(t, 13, len(tenant1Ue2DevTestTestComponentOverrideComponentDeps))
	assert.Equal(t, "catalog/terraform/services/service-1", tenant1Ue2DevTestTestComponentOverrideComponentDeps[0])
	assert.Equal(t, "catalog/terraform/services/service-1-override", tenant1Ue2DevTestTestComponentOverrideComponentDeps[1])
	assert.Equal(t, "catalog/terraform/services/service-2", tenant1Ue2DevTestTestComponentOverrideComponentDeps[2])
	assert.Equal(t, "catalog/terraform/services/service-2-override", tenant1Ue2DevTestTestComponentOverrideComponentDeps[3])
	assert.Equal(t, "catalog/terraform/spacelift-and-backend-override-1", tenant1Ue2DevTestTestComponentOverrideComponentDeps[4])
	assert.Equal(t, "catalog/terraform/tenant1-ue2-dev", tenant1Ue2DevTestTestComponentOverrideComponentDeps[5])
	assert.Equal(t, "catalog/terraform/test-component", tenant1Ue2DevTestTestComponentOverrideComponentDeps[6])
	assert.Equal(t, "catalog/terraform/test-component-override", tenant1Ue2DevTestTestComponentOverrideComponentDeps[7])
	assert.Equal(t, "catalog/terraform/test-component-override-3", tenant1Ue2DevTestTestComponentOverrideComponentDeps[8])
	assert.Equal(t, "globals/globals", tenant1Ue2DevTestTestComponentOverrideComponentDeps[9])
	assert.Equal(t, "globals/tenant1-globals", tenant1Ue2DevTestTestComponentOverrideComponentDeps[10])
	assert.Equal(t, "globals/ue2-globals", tenant1Ue2DevTestTestComponentOverrideComponentDeps[11])
	assert.Equal(t, "tenant1/ue2/dev", tenant1Ue2DevTestTestComponentOverrideComponentDeps[12])
	assert.Equal(t, 39, len(tenant1Ue2DevTestTestComponentOverrideComponentLabels))
	assert.Equal(t, "deps:stacks/catalog/terraform/services/service-1-override.yaml", tenant1Ue2DevTestTestComponentOverrideComponentLabels[25])
	assert.Equal(t, "deps:stacks/catalog/terraform/services/service-2.yaml", tenant1Ue2DevTestTestComponentOverrideComponentLabels[26])
	assert.Equal(t, "deps:stacks/catalog/terraform/services/service-2-override.yaml", tenant1Ue2DevTestTestComponentOverrideComponentLabels[27])
	assert.Equal(t, "deps:stacks/catalog/terraform/spacelift-and-backend-override-1.yaml", tenant1Ue2DevTestTestComponentOverrideComponentLabels[28])
	assert.Equal(t, "deps:stacks/catalog/terraform/tenant1-ue2-dev.yaml", tenant1Ue2DevTestTestComponentOverrideComponentLabels[29])
	assert.Equal(t, "deps:stacks/catalog/terraform/test-component.yaml", tenant1Ue2DevTestTestComponentOverrideComponentLabels[30])
	assert.Equal(t, "deps:stacks/catalog/terraform/test-component-override.yaml", tenant1Ue2DevTestTestComponentOverrideComponentLabels[31])
	assert.Equal(t, "deps:stacks/catalog/terraform/test-component-override-3.yaml", tenant1Ue2DevTestTestComponentOverrideComponentLabels[32])
	assert.Equal(t, "deps:stacks/globals/globals.yaml", tenant1Ue2DevTestTestComponentOverrideComponentLabels[33])
	assert.Equal(t, "deps:stacks/globals/tenant1-globals.yaml", tenant1Ue2DevTestTestComponentOverrideComponentLabels[34])
	assert.Equal(t, "deps:stacks/globals/ue2-globals.yaml", tenant1Ue2DevTestTestComponentOverrideComponentLabels[35])
	assert.Equal(t, "deps:stacks/tenant1/ue2/dev.yaml", tenant1Ue2DevTestTestComponentOverrideComponentLabels[36])
	assert.Equal(t, "folder:component/test/test-component-override", tenant1Ue2DevTestTestComponentOverrideComponentLabels[37])
	assert.Equal(t, "folder:tenant1/ue2/dev", tenant1Ue2DevTestTestComponentOverrideComponentLabels[38])
	assert.Equal(t, "test-component-override-workspace-override", tenant1Ue2DevTestTestComponentOverrideTerraformWorkspace)

	newTenant1Ue2DevTestTestComponentOverrideComponent2 := spaceliftStacks["tenant1-ue2-dev-new-component"].(map[string]interface{})
	newTenant1Ue2DevTestTestComponentOverrideComponent2InfrastructureStackName := newTenant1Ue2DevTestTestComponentOverrideComponent2["stack"].(string)
	assert.Equal(t, "tenant1-ue2-dev", newTenant1Ue2DevTestTestComponentOverrideComponent2InfrastructureStackName)

	yamlSpaceliftStacks, err := yaml.Marshal(spaceliftStacks)
	assert.Nil(t, err)
	t.Log(string(yamlSpaceliftStacks))
}

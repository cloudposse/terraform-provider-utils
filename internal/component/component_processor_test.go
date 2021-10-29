package component

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestComponentProcessor(t *testing.T) {
	component := "test/test-component-override"
	stack := "tenant1-ue2-dev"

	var tenant1Ue2DevTestTestComponentOverrideComponent, err = ProcessComponent(component, stack)
	assert.Nil(t, err)
	tenant1Ue2DevTestTestComponentOverrideComponentBackend := tenant1Ue2DevTestTestComponentOverrideComponent["backend"].(map[interface{}]interface{})
	tenant1Ue2DevTestTestComponentOverrideComponentBaseComponent := tenant1Ue2DevTestTestComponentOverrideComponent["component"].(string)
	tenant1Ue2DevTestTestComponentOverrideComponentBackendWorkspaceKeyPrefix := tenant1Ue2DevTestTestComponentOverrideComponentBackend["workspace_key_prefix"].(string)
	tenant1Ue2DevTestTestComponentOverrideComponentDeps := tenant1Ue2DevTestTestComponentOverrideComponent["deps"].([]string)
	assert.Equal(t, "test-test-component", tenant1Ue2DevTestTestComponentOverrideComponentBackendWorkspaceKeyPrefix)
	assert.Equal(t, "test/test-component", tenant1Ue2DevTestTestComponentOverrideComponentBaseComponent)
	assert.Equal(t, 11, len(tenant1Ue2DevTestTestComponentOverrideComponentDeps))
	assert.Equal(t, "catalog/terraform/services/service-1", tenant1Ue2DevTestTestComponentOverrideComponentDeps[0])
	assert.Equal(t, "catalog/terraform/services/service-1-override", tenant1Ue2DevTestTestComponentOverrideComponentDeps[1])
	assert.Equal(t, "catalog/terraform/services/service-2", tenant1Ue2DevTestTestComponentOverrideComponentDeps[2])
	assert.Equal(t, "catalog/terraform/services/service-2-override", tenant1Ue2DevTestTestComponentOverrideComponentDeps[3])
	assert.Equal(t, "catalog/terraform/tenant1-ue2-dev", tenant1Ue2DevTestTestComponentOverrideComponentDeps[4])
	assert.Equal(t, "catalog/terraform/test-component", tenant1Ue2DevTestTestComponentOverrideComponentDeps[5])
	assert.Equal(t, "catalog/terraform/test-component-override", tenant1Ue2DevTestTestComponentOverrideComponentDeps[6])
	assert.Equal(t, "globals/globals", tenant1Ue2DevTestTestComponentOverrideComponentDeps[7])
	assert.Equal(t, "globals/tenant1-globals", tenant1Ue2DevTestTestComponentOverrideComponentDeps[8])
	assert.Equal(t, "globals/ue2-globals", tenant1Ue2DevTestTestComponentOverrideComponentDeps[9])
	assert.Equal(t, "tenant1/ue2/dev", tenant1Ue2DevTestTestComponentOverrideComponentDeps[10])

	yamlConfig, err := yaml.Marshal(tenant1Ue2DevTestTestComponentOverrideComponent)
	assert.Nil(t, err)
	t.Log(string(yamlConfig))
}

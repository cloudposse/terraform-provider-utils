package component

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestComponentProcessor(t *testing.T) {
	component := "test/test-component-override"
	stack := "tenant1-ue2-dev"

	var componentConfig, err = ProcessComponent(component, stack)
	assert.Nil(t, err)

	yamlConfig, err := yaml.Marshal(componentConfig)
	assert.Nil(t, err)
	t.Log(string(yamlConfig))
}

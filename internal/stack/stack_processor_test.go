package stack

import (
	c "github.com/cloudposse/terraform-provider-utils/internal/convert"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestStackProcessor(t *testing.T) {
	filePaths := []string{
		"../../examples/data-sources/utils_stack_config_yaml/stacks/uw2-test.yaml",
	}

	yamlResult, err := ProcessYAMLConfigFiles(filePaths)
	assert.Nil(t, err)
	assert.Equal(t, len(yamlResult), 1)

	mapResult, err := c.YAMLToMapOfInterfaces(yamlResult[0])
	assert.Nil(t, err)

	yamlConfig, err := yaml.Marshal(mapResult)
	assert.Nil(t, err)

	t.Log(string(yamlConfig))
}

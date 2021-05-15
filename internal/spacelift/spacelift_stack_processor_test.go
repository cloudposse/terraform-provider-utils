package spacelift

import (
	u "github.com/cloudposse/terraform-provider-utils/internal/utils"
	"gopkg.in/yaml.v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpaceliftStackProcessor(t *testing.T) {
	filePaths := []string{
		"../../examples/data-sources/utils_stack_config_yaml/stacks/uw2-dev.yaml",
		"../../examples/data-sources/utils_stack_config_yaml/stacks/uw2-prod.yaml",
		"../../examples/data-sources/utils_stack_config_yaml/stacks/uw2-staging.yaml",
		"../../examples/data-sources/utils_stack_config_yaml/stacks/uw2-uat.yaml",
	}

	processStackDeps := true
	processComponentDeps := true

	var mapResult, err = ProcessSpaceliftConfigFiles(filePaths, processStackDeps, processComponentDeps)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(mapResult))

	mapResultKeys := u.StringKeysFromMap(mapResult)
	assert.Equal(t, "uw2-dev", mapResultKeys[0])
	assert.Equal(t, "uw2-prod", mapResultKeys[1])
	assert.Equal(t, "uw2-staging", mapResultKeys[2])
	assert.Equal(t, "uw2-uat", mapResultKeys[3])

	yamlConfig, err := yaml.Marshal(mapResult)
	assert.Nil(t, err)
	t.Log(string(yamlConfig))
}

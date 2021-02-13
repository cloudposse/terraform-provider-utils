package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackProcessor(t *testing.T) {
	filePaths := []string{
		"../../examples/data-sources/utils_stack_config_yaml/stacks/uw2-test.yaml",
	}

	result, err := ProcessYAMLConfigFiles(filePaths)
	assert.Nil(t, err)
	t.Log(result)
}

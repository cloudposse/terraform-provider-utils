package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYAMLToMap(t *testing.T) {
	input := `---
hello: world`
	result, err := YAMLToMap(input)
	assert.Nil(t, err)
	assert.Equal(t, result["hello"], "world")
}

func TestYAMLToMapRedPath(t *testing.T) {
	input := "Not YAML"
	_, err := YAMLToMap(input)
	assert.NotNil(t, err)
}

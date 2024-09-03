package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"

	u "github.com/cloudposse/atmos/pkg/utils"
)

func TestYAMLToMapOfInterfaces(t *testing.T) {
	input := `---
hello: world`
	result, err := u.UnmarshalYAML[map[any]any](input)
	assert.Nil(t, err)
	assert.Equal(t, result["hello"], "world")
}

func TestYAMLToMapOfInterfacesRedPath(t *testing.T) {
	input := "Not YAML"
	_, err := u.UnmarshalYAML[map[any]any](input)
	assert.NotNil(t, err)
}

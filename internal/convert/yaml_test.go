package convert

import (
	c "github.com/cloudposse/terraform-provider-utils/pkg/convert"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYAMLToMapOfInterfaces(t *testing.T) {
	input := `---
hello: world`
	result, err := c.YAMLToMapOfInterfaces(input)
	assert.Nil(t, err)
	assert.Equal(t, result["hello"], "world")
}

func TestYAMLToMapOfInterfacesRedPath(t *testing.T) {
	input := "Not YAML"
	_, err := c.YAMLToMapOfInterfaces(input)
	assert.NotNil(t, err)
}

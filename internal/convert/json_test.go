package convert

import (
	c "github.com/cloudposse/terraform-provider-utils/pkg/convert"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONToMapOfInterfaces(t *testing.T) {
	input := "{\"hello\": \"world\"}"
	result, err := c.JSONToMapOfInterfaces(input)
	assert.Nil(t, err)
	assert.Equal(t, result["hello"], "world")
}

func TestJSONToMapOfInterfacesRedPath(t *testing.T) {
	input := "Not JSON"
	_, err := c.JSONToMapOfInterfaces(input)
	assert.NotNil(t, err)
}

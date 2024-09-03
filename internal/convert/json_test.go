package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"

	u "github.com/cloudposse/atmos/pkg/utils"
)

func TestJSONToMapOfInterfaces(t *testing.T) {
	input := "{\"hello\": \"world\"}"
	result, err := u.JSONToMapOfInterfaces(input)
	assert.Nil(t, err)
	assert.Equal(t, result["hello"], "world")
}

func TestJSONToMapOfInterfacesRedPath(t *testing.T) {
	input := "Not JSON"
	_, err := u.JSONToMapOfInterfaces(input)
	assert.NotNil(t, err)
}

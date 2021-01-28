package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONToMap(t *testing.T) {
	input := "{\"hello\": \"world\"}"
	result, err := JSONToMap(input)
	assert.Nil(t, err)
	assert.Equal(t, result["hello"], "world")
}

func TestJSONToMapRedPath(t *testing.T) {
	input := "Not JSON"
	_, err := JSONToMap(input)
	assert.NotNil(t, err)
}

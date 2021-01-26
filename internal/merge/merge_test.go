package merge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeBasic(t *testing.T) {
	map1 := map[string]interface{}{"foo": "bar"}
	map2 := map[string]interface{}{"baz": "bat"}

	inputs := []interface{}{map1, map2}
	expected := map[string]interface{}{"foo": "bar", "baz": "bat"}

	result, err := Merge(inputs)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestMergeBasicOverride(t *testing.T) {
	map1 := map[string]interface{}{"foo": "bar"}
	map2 := map[string]interface{}{"baz": "bat"}
	map3 := map[string]interface{}{"foo": "ood"}

	inputs := []interface{}{map1, map2, map3}
	expected := map[string]interface{}{"foo": "ood", "baz": "bat"}

	result, err := Merge(inputs)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

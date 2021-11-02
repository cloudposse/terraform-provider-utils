package merge

import (
	m "github.com/cloudposse/atmos/pkg/merge"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeBasic(t *testing.T) {
	map1 := map[interface{}]interface{}{"foo": "bar"}
	map2 := map[interface{}]interface{}{"baz": "bat"}

	inputs := []map[interface{}]interface{}{map1, map2}
	expected := map[interface{}]interface{}{"foo": "bar", "baz": "bat"}

	result, err := m.Merge(inputs)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestMergeBasicOverride(t *testing.T) {
	map1 := map[interface{}]interface{}{"foo": "bar"}
	map2 := map[interface{}]interface{}{"baz": "bat"}
	map3 := map[interface{}]interface{}{"foo": "ood"}

	inputs := []map[interface{}]interface{}{map1, map2, map3}
	expected := map[interface{}]interface{}{"foo": "ood", "baz": "bat"}

	result, err := m.Merge(inputs)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

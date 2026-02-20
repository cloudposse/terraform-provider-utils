package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceOfInterfacesToSliceOfStrings(t *testing.T) {
	input := []any{"a", "b", "c"}

	result, err := SliceOfInterfacesToSliceOfStrings(input)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(result))
	assert.Equal(t, "a", result[0])
	assert.Equal(t, "b", result[1])
	assert.Equal(t, "c", result[2])
}

func TestSliceOfInterfacesToSliceOfStringsEmpty(t *testing.T) {
	input := []any{}

	result, err := SliceOfInterfacesToSliceOfStrings(input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}

func TestSliceOfInterfacesToSliceOfStringsNil(t *testing.T) {
	_, err := SliceOfInterfacesToSliceOfStrings(nil)
	assert.NotNil(t, err)
	assert.Equal(t, "input must not be nil", err.Error())
}

func TestSliceOfInterfacesToSliceOfStringsNonString(t *testing.T) {
	input := []any{"a", 42, "c"}

	_, err := SliceOfInterfacesToSliceOfStrings(input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "element at index 1 is not a string")
}

func TestSetEnvEmpty(t *testing.T) {
	err := setEnv(map[string]any{})
	assert.Nil(t, err)
}

func TestSetEnvNil(t *testing.T) {
	err := setEnv(nil)
	assert.Nil(t, err)
}

func TestYAMLSliceOfInterfaceToSliceOfMapsBasic(t *testing.T) {
	input := []any{
		"key1: value1\nkey2: value2",
		"key3: value3",
	}

	result, err := YAMLSliceOfInterfaceToSliceOfMaps(input)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "value1", result[0]["key1"])
	assert.Equal(t, "value2", result[0]["key2"])
	assert.Equal(t, "value3", result[1]["key3"])
}

func TestYAMLSliceOfInterfaceToSliceOfMapsEmpty(t *testing.T) {
	input := []any{}

	result, err := YAMLSliceOfInterfaceToSliceOfMaps(input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}

func TestYAMLSliceOfInterfaceToSliceOfMapsSkipsNonStrings(t *testing.T) {
	input := []any{
		"key1: value1",
		42,
		"key2: value2",
	}

	result, err := YAMLSliceOfInterfaceToSliceOfMaps(input)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "value1", result[0]["key1"])
	assert.Equal(t, "value2", result[1]["key2"])
}

func TestYAMLSliceOfInterfaceToSliceOfMapsInvalidYAML(t *testing.T) {
	input := []any{
		"invalid: [yaml: broken",
	}

	_, err := YAMLSliceOfInterfaceToSliceOfMaps(input)
	assert.NotNil(t, err)
}

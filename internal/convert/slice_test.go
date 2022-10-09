package convert

import (
	c "github.com/cloudposse/terraform-provider-utils/pkg/convert"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceOfInterfacesToSliceOfStrings(t *testing.T) {
	var input []any
	input = append(input, "a")
	input = append(input, "b")
	input = append(input, "c")

	result, err := c.SliceOfInterfacesToSliceOfStrings(input)

	assert.Nil(t, err)
	assert.Equal(t, len(input), len(result))
	assert.Equal(t, input[0].(string), result[0])
	assert.Equal(t, input[1].(string), result[1])
	assert.Equal(t, input[2].(string), result[2])
}

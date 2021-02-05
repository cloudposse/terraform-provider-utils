package convert

import "github.com/pkg/errors"

// SliceOfInterfacesToSliceOfStrings takes a slice of interfaces and converts it to a slice of strings
func SliceOfInterfacesToSliceOfStrings(input []interface{}) ([]string, error) {
	if input == nil {
		return nil, errors.New("imput must not be nil")
	}

	output := make([]string, 0)
	for _, current := range input {
		output = append(output, current.(string))
	}

	return output, nil
}

package provider

import (
	"os"

	"github.com/pkg/errors"
)

func setEnv(envMap map[string]any) error {
	for k, v := range envMap {
		val := v.(string)
		err := os.Setenv(k, val)
		if err != nil {
			return err
		}
	}
	return nil
}

// SliceOfInterfacesToSliceOfStrings takes a slice of interfaces and converts it to a slice of strings.
// Returns an error if the input is nil or contains non-string elements.
func SliceOfInterfacesToSliceOfStrings(input []any) ([]string, error) {
	if input == nil {
		return nil, errors.New("input must not be nil")
	}

	output := make([]string, 0)
	for i, current := range input {
		str, ok := current.(string)
		if !ok {
			return nil, errors.Errorf("element at index %d is not a string: %T", i, current)
		}
		output = append(output, str)
	}

	return output, nil
}

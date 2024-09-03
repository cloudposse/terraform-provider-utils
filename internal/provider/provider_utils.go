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

// SliceOfInterfacesToSliceOfStrings takes a slice of interfaces and converts it to a slice of strings
func SliceOfInterfacesToSliceOfStrings(input []any) ([]string, error) {
	if input == nil {
		return nil, errors.New("input must not be nil")
	}

	output := make([]string, 0)
	for _, current := range input {
		output = append(output, current.(string))
	}

	return output, nil
}

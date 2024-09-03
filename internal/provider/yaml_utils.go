package provider

import (
	u "github.com/cloudposse/atmos/pkg/utils"
)

// YAMLSliceOfInterfaceToSliceOfMaps takes a slice of interfaces as input and returns a slice of map[any]any
func YAMLSliceOfInterfaceToSliceOfMaps(input []any) ([]map[string]any, error) {
	output := make([]map[string]any, 0)
	for _, current := range input {
		// Apply YAMLToMap only if string is passed
		if currentYaml, ok := current.(string); ok {
			data, err := u.UnmarshalYAML[map[string]any](currentYaml)
			if err != nil {
				return nil, err
			}
			output = append(output, data)
		}
	}
	return output, nil
}

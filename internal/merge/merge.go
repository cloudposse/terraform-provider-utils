package merge

import (
	"github.com/imdario/mergo"
)

// Merge takes a list of maps as input and returns a single map with the merged contents
func Merge(inputs []map[string]interface{}) (map[string]interface{}, error) {
	merged := map[string]interface{}{}

	for index := range inputs {
		if err := mergo.Merge(&merged, inputs[index], mergo.WithOverride, mergo.WithOverwriteWithEmptyValue); err != nil {
			return nil, err
		}
	}

	return merged, nil
}

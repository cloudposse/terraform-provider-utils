package merge

import (
	"github.com/imdario/mergo"
)

// Merge takes a list of maps as input and returns a single map with the merged contents
func Merge(inputs []map[interface{}]interface{}) (map[interface{}]interface{}, error) {
	merged := map[interface{}]interface{}{}

	for index := range inputs {
		if err := mergo.Merge(&merged, inputs[index], mergo.WithOverride, mergo.WithOverwriteWithEmptyValue, mergo.WithTypeCheck); err != nil {
			return nil, err
		}
	}

	return merged, nil
}

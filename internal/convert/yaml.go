package convert

import "gopkg.in/yaml.v2"

// YAMLToMap takes a YAML string as input and returns a map[string]interface{}
func YAMLToMap(input string) (map[string]interface{}, error) {
	var data map[string]interface{}
	byt := []byte(input)

	if err := yaml.Unmarshal(byt, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// YAMLSliceOfInterfaceToSliceOfMaps takes a slice of JSON strings as input and returns a slice of map[string]interface{}
func YAMLSliceOfInterfaceToSliceOfMaps(input []interface{}) ([]map[string]interface{}, error) {
	outputMap := make([]map[string]interface{}, 0)
	for _, current := range input {
		data, err := YAMLToMap(current.(string))
		if err != nil {
			return nil, err
		}
		outputMap = append(outputMap, data)
	}
	return outputMap, nil
}

package provider

import "encoding/json"

// JSONToMapOfInterfaces takes a JSON string as input and returns a map[string]any
func JSONToMapOfInterfaces(input string) (map[string]any, error) {
	var data map[string]any
	byt := []byte(input)

	if err := json.Unmarshal(byt, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// JSONSliceOfInterfaceToSliceOfMaps takes a slice of JSON strings as input and returns a slice of map[any]any
func JSONSliceOfInterfaceToSliceOfMaps(input []any) ([]map[string]any, error) {
	outputMap := make([]map[string]any, 0)
	for _, current := range input {
		data, err := JSONToMapOfInterfaces(current.(string))
		if err != nil {
			return nil, err
		}

		map2 := map[string]any{}

		for k, v := range data {
			map2[k] = v
		}

		outputMap = append(outputMap, map2)
	}
	return outputMap, nil
}

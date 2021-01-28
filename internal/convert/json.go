package convert

import "encoding/json"

// JSONToMap takes a JSON string as input and returns a map[string]interface{}
func JSONToMap(input string) (map[string]interface{}, error) {
	var data map[string]interface{}
	byt := []byte(input)

	if err := json.Unmarshal(byt, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// JSONSliceOfInterfaceToSliceOfMaps takes a slice of JSON strings as input and returns a slice of map[string]interface{}
func JSONSliceOfInterfaceToSliceOfMaps(input []interface{}) ([]map[string]interface{}, error) {
	outputMap := make([]map[string]interface{}, 0)
	for _, current := range input {
		data, err := JSONToMap(current.(string))
		if err != nil {
			return nil, err
		}
		outputMap = append(outputMap, data)
	}
	return outputMap, nil
}

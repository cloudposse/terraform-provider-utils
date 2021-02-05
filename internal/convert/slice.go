package convert

// SliceOfInterfacesToSliceOfStrings takes a slice of interfaces and converts it to a slice of strings
func SliceOfInterfacesToSliceOfStrings(input []interface{}) ([]string, error) {
	output := make([]string, 0)
	for _, current := range input {
		output = append(output, current.(string))
	}
	return output, nil
}

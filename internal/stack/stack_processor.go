package stack

import (
	c "github.com/cloudposse/terraform-provider-utils/internal/convert"
	"io/ioutil"
)

// ProcessYAMLConfigFile takes a path to a YAML config file, processes and deep-merges all imports,
// and returns stack config as map[string]interface{}
func ProcessYAMLConfigFile(path string) (map[string]interface{}, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	configMap, err := c.YAMLToMap(string(content))
	if err != nil {
		return nil, err
	}

	return configMap, nil
}

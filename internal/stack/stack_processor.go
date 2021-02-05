package stack

import (
	c "github.com/cloudposse/terraform-provider-utils/internal/convert"
	m "github.com/cloudposse/terraform-provider-utils/internal/merge"
	"io/ioutil"
	"path"
)

// ProcessYAMLConfigFile takes a path to a YAML config file, processes and deep-merges all imports,
// and returns stack config as map[string]interface{}
func ProcessYAMLConfigFile(filePath string) (map[string]interface{}, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	config, err := c.YAMLToMap(string(content))
	if err != nil {
		return nil, err
	}

	var configs []map[string]interface{}
	dir := path.Dir(filePath)

	// Find and process all imports
	if imports, ok := config["import"]; ok {
		for _, i := range imports.([]interface{}) {
			p := path.Join(dir, i.(string)+".yaml")

			yamlConfig, err := ProcessYAMLConfigFile(p)
			if err != nil {
				return nil, err
			}
			configs = append(configs, yamlConfig)
		}
	}

	configs = append(configs, config)

	// Deep-merge the config file and all the imports
	result, err := m.Merge(configs)
	if err != nil {
		return nil, err
	}

	return result, nil
}

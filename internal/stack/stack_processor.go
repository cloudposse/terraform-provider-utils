package stack

import (
	c "github.com/cloudposse/terraform-provider-utils/internal/convert"
	m "github.com/cloudposse/terraform-provider-utils/internal/merge"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
)

// ProcessYAMLConfigFiles takes a list of paths to YAML config files, processes and deep-merges all imports,
// and returns a list of stack configs
func ProcessYAMLConfigFiles(filePaths []string) ([]string, error) {
	var result []string

	for _, p := range filePaths {
		config, err := ProcessYAMLConfigFile(p)
		if err != nil {
			return nil, err
		}

		finalConfig, err := ProcessConfig(config)
		if err != nil {
			return nil, err
		}

		yamlConfig, err := yaml.Marshal(finalConfig)
		if err != nil {
			return nil, err
		}

		result = append(result, string(yamlConfig))
	}

	return result, nil
}

// ProcessYAMLConfigFile takes a path to a YAML config file,
// recursively processes and deep-merges all imports,
// and returns stack config as map[interface{}]interface{}
func ProcessYAMLConfigFile(filePath string) (map[interface{}]interface{}, error) {
	var configs []map[interface{}]interface{}
	dir := path.Dir(filePath)

	stackYamlConfig, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	stackMapConfig, err := c.YAMLToMapOfInterfaces(string(stackYamlConfig))
	if err != nil {
		return nil, err
	}

	// Find and process all imports
	if imports, ok := stackMapConfig["import"]; ok {
		for _, i := range imports.([]interface{}) {
			p := path.Join(dir, i.(string)+".yaml")

			yamlConfig, err := ProcessYAMLConfigFile(p)
			if err != nil {
				return nil, err
			}
			configs = append(configs, yamlConfig)
		}
	}

	configs = append(configs, stackMapConfig)

	// Deep-merge the config file and the imports
	result, err := m.Merge(configs)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ProcessConfig takes a raw stack config, deep-merges all variables and backends,
// and returns the final stack configuration for all Terraform and helmfile components
func ProcessConfig(config map[interface{}]interface{}) (map[interface{}]interface{}, error) {
	globalVars := map[interface{}]interface{}{}
	terraformVars := map[interface{}]interface{}{}
	helmfileVars := map[interface{}]interface{}{}
	backendType := "s3"
	backend := map[interface{}]interface{}{}
	terraformComponents := map[interface{}]interface{}{}
	helmfileComponents := map[interface{}]interface{}{}
	allComponents := map[interface{}]interface{}{}

	if i, ok := config["vars"]; ok {
		globalVars = i.(map[interface{}]interface{})
	}

	if i, ok := config["terraform"].(map[interface{}]interface{})["vars"]; ok {
		terraformVars = i.(map[interface{}]interface{})
	}

	if i, ok := config["helmfile"].(map[interface{}]interface{})["vars"]; ok {
		helmfileVars = i.(map[interface{}]interface{})
	}

	if i, ok := config["terraform"].(map[interface{}]interface{})["backend_type"]; ok {
		backendType = i.(string)
	}

	if i, ok := config["terraform"].(map[interface{}]interface{})["backend"].(map[interface{}]interface{})[backendType]; ok {
		backend = i.(map[interface{}]interface{})
	}

	if i, ok := config["components"].(map[interface{}]interface{})["terraform"].(map[interface{}]interface{}); ok {
		for k, v := range i {
			componentVars := map[interface{}]interface{}{}
			if i2, ok2 := v.(map[interface{}]interface{})["vars"]; ok2 {
				componentVars = i2.(map[interface{}]interface{})
			}

			componentBackend := map[interface{}]interface{}{}
			if i2, ok2 := v.(map[interface{}]interface{})["backend"].(map[interface{}]interface{})[backendType]; ok2 {
				componentBackend = i2.(map[interface{}]interface{})
			}

			allComponentVars, err := m.Merge([]map[interface{}]interface{}{globalVars, terraformVars, componentVars})
			if err != nil {
				return nil, err
			}

			allComponentBackend, err := m.Merge([]map[interface{}]interface{}{backend, componentBackend})
			if err != nil {
				return nil, err
			}

			comp := map[interface{}]interface{}{}
			comp["vars"] = allComponentVars
			comp["backend_type"] = backendType
			comp["backend"] = allComponentBackend
			terraformComponents[k] = comp
		}
	}

	if i, ok := config["components"].(map[interface{}]interface{})["helmfile"].(map[interface{}]interface{}); ok {
		for k, v := range i {
			componentVars := map[interface{}]interface{}{}
			if i2, ok2 := v.(map[interface{}]interface{})["vars"]; ok2 {
				componentVars = i2.(map[interface{}]interface{})
			}

			allComponentVars, err := m.Merge([]map[interface{}]interface{}{globalVars, helmfileVars, componentVars})
			if err != nil {
				return nil, err
			}

			comp := map[interface{}]interface{}{}
			comp["vars"] = allComponentVars
			helmfileComponents[k] = comp
		}
	}

	allComponents["terraform"] = terraformComponents
	allComponents["helmfile"] = helmfileComponents

	result := map[interface{}]interface{}{
		"config":     config,
		"components": allComponents,
	}

	return result, nil
}

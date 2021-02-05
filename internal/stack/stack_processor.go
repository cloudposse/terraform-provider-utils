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
// and returns stack config as map[string]interface{}
func ProcessYAMLConfigFile(filePath string) (map[string]interface{}, error) {
	var configs []map[string]interface{}
	dir := path.Dir(filePath)

	stackYamlConfig, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	stackMapConfig, err := c.YAMLToMap(string(stackYamlConfig))
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
func ProcessConfig(config map[string]interface{}) (map[string]interface{}, error) {
	globalVars := map[string]interface{}{}
	terraformVars := map[string]interface{}{}
	helmfileVars := map[string]interface{}{}
	backendType := "s3"
	backend := map[string]interface{}{}
	terraformComponents := map[string]interface{}{}
	helmfileComponents := map[string]interface{}{}
	allComponents := map[string]interface{}{}

	if i, ok := config["vars"]; ok {
		globalVars = i.(map[string]interface{})
	}

	if i, ok := config["terraform"].(map[string]interface{})["vars"]; ok {
		terraformVars = i.(map[string]interface{})
	}

	if i, ok := config["helmfile"].(map[string]interface{})["vars"]; ok {
		helmfileVars = i.(map[string]interface{})
	}

	if i, ok := config["terraform"].(map[string]interface{})["backend_type"]; ok {
		backendType = i.(string)
	}

	if i, ok := config["terraform"].(map[string]interface{})["backend"].(map[string]interface{})[backendType]; ok {
		backend = i.(map[string]interface{})
	}

	if i, ok := config["components"].(map[string]interface{})["terraform"].(map[string]interface{}); ok {
		for k, v := range i {
			componentVars := map[string]interface{}{}
			if i2, ok2 := v.(map[string]interface{})["vars"]; ok2 {
				componentVars = i2.(map[string]interface{})
			}

			componentBackend := map[string]interface{}{}
			if i2, ok2 := v.(map[string]interface{})["backend"].(map[string]interface{})[backendType]; ok2 {
				componentBackend = i2.(map[string]interface{})
			}

			allComponentVars, err := m.Merge([]map[string]interface{}{globalVars, terraformVars, componentVars})
			if err != nil {
				return nil, err
			}

			allComponentBackend, err := m.Merge([]map[string]interface{}{backend, componentBackend})
			if err != nil {
				return nil, err
			}

			comp := map[string]interface{}{}
			comp["vars"] = allComponentVars
			comp["backend_type"] = backendType
			comp["backend"] = allComponentBackend
			terraformComponents[k] = comp
		}
	}

	if i, ok := config["components"].(map[string]interface{})["helmfile"].(map[string]interface{}); ok {
		for k, v := range i {
			componentVars := map[string]interface{}{}
			if i2, ok2 := v.(map[string]interface{})["vars"]; ok2 {
				componentVars = i2.(map[string]interface{})
			}

			allComponentVars, err := m.Merge([]map[string]interface{}{globalVars, helmfileVars, componentVars})
			if err != nil {
				return nil, err
			}

			comp := map[string]interface{}{}
			comp["vars"] = allComponentVars
			helmfileComponents[k] = comp
		}
	}

	allComponents["terraform"] = terraformComponents
	allComponents["helmfile"] = helmfileComponents

	result := map[string]interface{}{
		"config":     config,
		"components": allComponents,
	}

	return result, nil
}

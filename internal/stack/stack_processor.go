package stack

import (
	c "github.com/cloudposse/terraform-provider-utils/internal/convert"
	m "github.com/cloudposse/terraform-provider-utils/internal/merge"
	"github.com/pkg/errors"
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

		finalConfig, err := ProcessConfig(p, config)
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
func ProcessConfig(stack string, config map[interface{}]interface{}) (map[interface{}]interface{}, error) {
	globalVars := map[interface{}]interface{}{}
	globalSettings := map[interface{}]interface{}{}
	terraformVars := map[interface{}]interface{}{}
	terraformSettings := map[interface{}]interface{}{}
	helmfileVars := map[interface{}]interface{}{}
	helmfileSettings := map[interface{}]interface{}{}
	backendType := "s3"
	backend := map[interface{}]interface{}{}
	terraformComponents := map[string]interface{}{}
	helmfileComponents := map[string]interface{}{}
	allComponents := map[string]interface{}{}

	if i, ok := config["vars"]; ok {
		globalVars = i.(map[interface{}]interface{})
	}

	if i, ok := config["settings"]; ok {
		globalSettings = i.(map[interface{}]interface{})
	}

	if i, ok := config["terraform"].(map[interface{}]interface{})["vars"]; ok {
		terraformVars = i.(map[interface{}]interface{})
	}

	if i, ok := config["terraform"].(map[interface{}]interface{})["settings"]; ok {
		terraformSettings = i.(map[interface{}]interface{})
	}

	if i, ok := config["terraform"].(map[interface{}]interface{})["backend_type"]; ok {
		backendType = i.(string)
	}

	if i, ok := config["terraform"].(map[interface{}]interface{})["backend"]; ok {
		if backendSection, backendSectionExist := i.(map[interface{}]interface{})[backendType]; backendSectionExist {
			backend = backendSection.(map[interface{}]interface{})
		}
	}

	if i, ok := config["helmfile"].(map[interface{}]interface{})["vars"]; ok {
		helmfileVars = i.(map[interface{}]interface{})
	}

	if i, ok := config["helmfile"].(map[interface{}]interface{})["settings"]; ok {
		helmfileSettings = i.(map[interface{}]interface{})
	}

	if allTerraformComponents, ok := config["components"].(map[interface{}]interface{})["terraform"]; ok {
		allTerraformComponentsMap := allTerraformComponents.(map[interface{}]interface{})
		for component, v := range allTerraformComponentsMap {
			componentMap := v.(map[interface{}]interface{})

			componentVars := map[interface{}]interface{}{}
			if i, ok2 := componentMap["vars"]; ok2 {
				componentVars = i.(map[interface{}]interface{})
			}

			componentSettings := map[interface{}]interface{}{}
			if i, ok2 := componentMap["settings"]; ok2 {
				componentSettings = i.(map[interface{}]interface{})
			}

			componentBackend := map[interface{}]interface{}{}
			if i, ok2 := componentMap["backend"]; ok2 {
				componentBackend = i.(map[interface{}]interface{})[backendType].(map[interface{}]interface{})
			}

			baseComponentVars := map[interface{}]interface{}{}
			baseComponentBackend := map[interface{}]interface{}{}
			baseComponentName := ""

			if baseComponent, baseComponentExist := componentMap["component"]; baseComponentExist {
				baseComponentName = baseComponent.(string)

				if baseComponentSection, baseComponentSectionExist := allTerraformComponentsMap[baseComponentName]; baseComponentSectionExist {
					baseComponentMap := baseComponentSection.(map[interface{}]interface{})
					baseComponentVars = baseComponentMap["vars"].(map[interface{}]interface{})

					if baseComponentBackendSection, baseComponentBackendSectionExist := baseComponentMap["backend"]; baseComponentBackendSectionExist {
						baseComponentBackend = baseComponentBackendSection.(map[interface{}]interface{})[backendType].(map[interface{}]interface{})
					}
				} else {
					return nil, errors.New("Terraform component '" + component.(string) + "' defines attribute 'component: " +
						baseComponentName + "', " + "but `" + baseComponentName + "' is not defined in the stack '" + stack + "'")
				}
			}

			finalComponentVars, err := m.Merge([]map[interface{}]interface{}{globalVars, terraformVars, baseComponentVars, componentVars})
			if err != nil {
				return nil, err
			}

			finalComponentSettings, err := m.Merge([]map[interface{}]interface{}{globalSettings, terraformSettings, componentSettings})
			if err != nil {
				return nil, err
			}

			finalComponentBackend, err := m.Merge([]map[interface{}]interface{}{backend, baseComponentBackend, componentBackend})
			if err != nil {
				return nil, err
			}

			comp := map[string]interface{}{}
			comp["vars"] = finalComponentVars
			comp["settings"] = finalComponentSettings
			comp["backend_type"] = backendType
			comp["backend"] = finalComponentBackend

			if baseComponentName != "" {
				comp["component"] = baseComponentName
			}

			terraformComponents[component.(string)] = comp
		}
	}

	if allHelmfileComponents, ok := config["components"].(map[interface{}]interface{})["helmfile"]; ok {
		allHelmfileComponentsMap := allHelmfileComponents.(map[interface{}]interface{})
		for component, v := range allHelmfileComponentsMap {
			componentMap := v.(map[interface{}]interface{})

			componentVars := map[interface{}]interface{}{}
			if i2, ok2 := componentMap["vars"]; ok2 {
				componentVars = i2.(map[interface{}]interface{})
			}

			componentSettings := map[interface{}]interface{}{}
			if i, ok2 := componentMap["settings"]; ok2 {
				componentSettings = i.(map[interface{}]interface{})
			}

			finalComponentVars, err := m.Merge([]map[interface{}]interface{}{globalVars, helmfileVars, componentVars})
			if err != nil {
				return nil, err
			}

			finalComponentSettings, err := m.Merge([]map[interface{}]interface{}{globalSettings, helmfileSettings, componentSettings})
			if err != nil {
				return nil, err
			}

			comp := map[string]interface{}{}
			comp["vars"] = finalComponentVars
			comp["settings"] = finalComponentSettings
			helmfileComponents[component.(string)] = comp
		}
	}

	allComponents["terraform"] = terraformComponents
	allComponents["helmfile"] = helmfileComponents

	result := map[interface{}]interface{}{
		"components": allComponents,
	}

	return result, nil
}

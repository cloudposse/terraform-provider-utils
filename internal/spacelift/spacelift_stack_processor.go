package spacelift

import (
	"fmt"
	s "github.com/cloudposse/terraform-provider-utils/internal/stack"
)

// CreateSpaceliftStacks takes a list of paths to YAML config files, processes and deep-merges all imports,
// and returns a map of Spacelift stack configs
func CreateSpaceliftStacks(filePaths []string, processStackDeps bool, processComponentDeps bool) (map[string]interface{}, error) {
	var _, mapResult, err = s.ProcessYAMLConfigFiles(filePaths, processStackDeps, processComponentDeps)
	if err != nil {
		return nil, err
	}
	return TransformStackConfigToSpaceliftStacks(mapResult)
}

// TransformStackConfigToSpaceliftStacks takes a a map of stack configs and transforms it to a map of Spacelift stacks
func TransformStackConfigToSpaceliftStacks(stacks map[string]interface{}) (map[string]interface{}, error) {
	res := map[string]interface{}{}

	for stackName, stackConfig := range stacks {
		config := stackConfig.(map[interface{}]interface{})
		imports := []string{}

		if i, ok := config["imports"]; ok {
			imports = i.([]string)
		}

		if i, ok := config["components"]; ok {
			componentsSection := i.(map[string]interface{})

			if terraformComponents, ok := componentsSection["terraform"]; ok {
				terraformComponentsMap := terraformComponents.(map[string]interface{})

				for component, v := range terraformComponentsMap {
					componentMap := v.(map[string]interface{})

					componentVars := map[interface{}]interface{}{}
					if i, ok2 := componentMap["vars"]; ok2 {
						componentVars = i.(map[interface{}]interface{})
					}

					componentSettings := map[interface{}]interface{}{}
					if i, ok2 := componentMap["settings"]; ok2 {
						componentSettings = i.(map[interface{}]interface{})
					}

					componentEnv := map[interface{}]interface{}{}
					if i, ok2 := componentMap["env"]; ok2 {
						componentEnv = i.(map[interface{}]interface{})
					}

					componentDeps := []string{}
					if i, ok2 := componentMap["deps"]; ok2 {
						componentDeps = i.([]string)
					}

					componentStacks := []string{}
					if i, ok2 := componentMap["stacks"]; ok2 {
						componentStacks = i.([]string)
					}

					spaceliftConfig := map[string]interface{}{}
					spaceliftConfig["component"] = component
					spaceliftConfig["stack"] = stackName
					spaceliftConfig["imports"] = imports
					spaceliftConfig["vars"] = componentVars
					spaceliftConfig["settings"] = componentSettings
					spaceliftConfig["env"] = componentEnv
					spaceliftConfig["deps"] = componentDeps
					spaceliftConfig["stacks"] = componentStacks

					spaceliftWorkspaceEnabled := false
					if i, ok2 := componentSettings["spacelift"]; ok2 {
						spaceliftSettings := i.(map[interface{}]interface{})

						if i3, ok3 := spaceliftSettings["workspace_enabled"]; ok3 {
							spaceliftWorkspaceEnabled = i3.(bool)
						}
					}
					spaceliftConfig["enabled"] = spaceliftWorkspaceEnabled

					baseComponentName := ""
					if baseComponent, baseComponentExist := componentMap["component"]; baseComponentExist {
						baseComponentName = baseComponent.(string)
					}
					spaceliftConfig["base_component"] = baseComponentName

					backendTypeName := ""
					if backendType, backendTypeExist := componentMap["backend_type"]; backendTypeExist {
						backendTypeName = backendType.(string)
					}
					spaceliftConfig["backend_type"] = backendTypeName

					componentBackend := map[interface{}]interface{}{}
					if i, ok2 := componentMap["backend"]; ok2 {
						componentBackend = i.(map[interface{}]interface{})
					}
					spaceliftConfig["backend"] = componentBackend

					var workspace string
					if backendTypeName == "s3" && baseComponentName == "" {
						workspace = stackName
					} else {
						workspace = fmt.Sprintf("%s-%s", stackName, component)
					}
					spaceliftConfig["workspace"] = workspace

					spaceliftStackName := fmt.Sprintf("%s-%s", stackName, component)
					res[spaceliftStackName] = spaceliftConfig
				}
			}
		}
	}

	return res, nil
}

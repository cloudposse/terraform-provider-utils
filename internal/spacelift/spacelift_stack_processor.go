package spacelift

import (
	"fmt"
	s "github.com/cloudposse/terraform-provider-utils/internal/stack"
	"strings"
)

// CreateSpaceliftStacks takes a list of paths to YAML config files, processes and deep-merges all imports,
// and returns a map of Spacelift stack configs
func CreateSpaceliftStacks(
	filePaths []string,
	processStackDeps bool,
	processComponentDeps bool,
	processImports bool,
	stackConfigPathTemplate string) (map[string]interface{}, error) {
	var _, mapResult, err = s.ProcessYAMLConfigFiles(filePaths, processStackDeps, processComponentDeps)
	if err != nil {
		return nil, err
	}
	return TransformStackConfigToSpaceliftStacks(mapResult, stackConfigPathTemplate, processImports)
}

// TransformStackConfigToSpaceliftStacks takes a a map of stack configs and transforms it to a map of Spacelift stacks
func TransformStackConfigToSpaceliftStacks(
	stacks map[string]interface{},
	stackConfigPathTemplate string,
	processImports bool) (map[string]interface{}, error) {

	res := map[string]interface{}{}

	for stackName, stackConfig := range stacks {
		config := stackConfig.(map[interface{}]interface{})
		imports := []string{}

		if processImports == true {
			if i, ok := config["imports"]; ok {
				imports = i.([]string)
			}
		}

		if i, ok := config["components"]; ok {
			componentsSection := i.(map[string]interface{})

			if terraformComponents, ok := componentsSection["terraform"]; ok {
				terraformComponentsMap := terraformComponents.(map[string]interface{})

				for component, v := range terraformComponentsMap {
					componentMap := v.(map[string]interface{})

					componentSettings := map[interface{}]interface{}{}
					if i, ok2 := componentMap["settings"]; ok2 {
						componentSettings = i.(map[interface{}]interface{})
					}

					spaceliftSettings := map[interface{}]interface{}{}
					spaceliftWorkspaceEnabled := false

					if i, ok2 := componentSettings["spacelift"]; ok2 {
						spaceliftSettings = i.(map[interface{}]interface{})

						if i3, ok3 := spaceliftSettings["workspace_enabled"]; ok3 {
							spaceliftWorkspaceEnabled = i3.(bool)
						}
					}

					// If Spacelift workspace is disabled, don't include it, continue to the next component
					if spaceliftWorkspaceEnabled == false {
						continue
					}

					spaceliftExplicitLabels := []interface{}{}
					if i, ok2 := spaceliftSettings["labels"]; ok2 {
						spaceliftExplicitLabels = i.([]interface{})
					}

					spaceliftDependsOn := []interface{}{}
					if i, ok2 := spaceliftSettings["depends_on"]; ok2 {
						spaceliftDependsOn = i.([]interface{})
					}

					spaceliftConfig := map[string]interface{}{}
					spaceliftConfig["enabled"] = spaceliftWorkspaceEnabled

					componentVars := map[interface{}]interface{}{}
					if i, ok2 := componentMap["vars"]; ok2 {
						componentVars = i.(map[interface{}]interface{})
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

					spaceliftConfig["component"] = component
					spaceliftConfig["stack"] = stackName
					spaceliftConfig["imports"] = imports
					spaceliftConfig["vars"] = componentVars
					spaceliftConfig["settings"] = componentSettings
					spaceliftConfig["env"] = componentEnv
					spaceliftConfig["deps"] = componentDeps
					spaceliftConfig["stacks"] = componentStacks

					baseComponentName := ""
					if baseComponent, baseComponentExist := componentMap["component"]; baseComponentExist {
						baseComponentName = baseComponent.(string)
					}
					spaceliftConfig["base_component"] = baseComponentName

					// backend
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

					// workspace
					var workspace string
					if backendTypeName == "s3" && baseComponentName == "" {
						workspace = stackName
					} else {
						workspace = fmt.Sprintf("%s-%s", stackName, component)
					}
					spaceliftConfig["workspace"] = workspace

					// labels
					labels := []string{}
					for _, v := range imports {
						labels = append(labels, fmt.Sprintf("import:"+stackConfigPathTemplate, v))
					}
					for _, v := range componentStacks {
						labels = append(labels, fmt.Sprintf("stack:"+stackConfigPathTemplate, v))
					}
					for _, v := range componentDeps {
						labels = append(labels, fmt.Sprintf("deps:"+stackConfigPathTemplate, v))
					}
					for _, v := range spaceliftExplicitLabels {
						labels = append(labels, v.(string))
					}
					for _, v := range spaceliftDependsOn {
						labels = append(labels, fmt.Sprintf("depends-on:%s-%s", stackName, v))
					}
					labels = append(labels, fmt.Sprintf("folder:component/%s", component))
					// Split on the first `-` and get the two parts: environment and stage
					stackNameParts := strings.SplitN(stackName, "-", 2)
					stackNamePartsLen := len(stackNameParts)
					if stackNamePartsLen == 2 {
						labels = append(labels, fmt.Sprintf("folder:%s/%s", stackNameParts[0], stackNameParts[1]))
					}
					spaceliftConfig["labels"] = labels

					// Add Spacelift stack config to the final map
					spaceliftStackName := fmt.Sprintf("%s-%s", stackName, component)
					res[spaceliftStackName] = spaceliftConfig
				}
			}
		}
	}

	return res, nil
}

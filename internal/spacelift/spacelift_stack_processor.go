package spacelift

import (
	s "github.com/cloudposse/terraform-provider-utils/internal/stack"
)

// ProcessSpaceliftConfigFiles takes a list of paths to YAML config files, processes and deep-merges all imports,
// and returns a map of Spacelift stack configs
func ProcessSpaceliftConfigFiles(filePaths []string, processStackDeps bool, processComponentDeps bool) ([]string, map[string]interface{}, error) {
	var listResult, mapResult, err = s.ProcessYAMLConfigFiles(filePaths, processStackDeps, processComponentDeps)
	return listResult, mapResult, err
}

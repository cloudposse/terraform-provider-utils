package component

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c "github.com/cloudposse/atmos/pkg/describe"
	u "github.com/cloudposse/atmos/pkg/utils"
)

func TestComponentProcessor(t *testing.T) {
	var err error
	var component string
	var stack string
	var yamlConfig string

	var tenant1Ue2DevTestTestComponent map[string]any
	component = "test/test-component"
	stack = "tenant1-ue2-dev"
	tenant1Ue2DevTestTestComponent, err = c.ProcessComponentInStack(component, stack, "", "")
	assert.Nil(t, err)
	tenant1Ue2DevTestTestComponentBackend := tenant1Ue2DevTestTestComponent["backend"].(map[string]any)
	tenant1Ue2DevTestTestComponentRemoteStateBackend := tenant1Ue2DevTestTestComponent["remote_state_backend"].(map[string]any)
	tenant1Ue2DevTestTestComponentBaseComponent := tenant1Ue2DevTestTestComponent["component"]
	tenant1Ue2DevTestTestComponentTerraformWorkspace := tenant1Ue2DevTestTestComponent["workspace"].(string)
	tenant1Ue2DevTestTestComponentWorkspace := tenant1Ue2DevTestTestComponent["workspace"].(string)
	tenant1Ue2DevTestTestComponentBackendWorkspaceKeyPrefix := tenant1Ue2DevTestTestComponentBackend["workspace_key_prefix"].(string)
	tenant1Ue2DevTestTestComponentRemoteStateBackendWorkspaceKeyPrefix := tenant1Ue2DevTestTestComponentRemoteStateBackend["workspace_key_prefix"].(string)
	tenant1Ue2DevTestTestComponentDeps := tenant1Ue2DevTestTestComponent["deps"].([]any)
	assert.Equal(t, "test-test-component", tenant1Ue2DevTestTestComponentBackendWorkspaceKeyPrefix)
	assert.Equal(t, "test-test-component", tenant1Ue2DevTestTestComponentRemoteStateBackendWorkspaceKeyPrefix)
	assert.Equal(t, "test/test-component", tenant1Ue2DevTestTestComponentBaseComponent)
	assert.Equal(t, "tenant1-ue2-dev", tenant1Ue2DevTestTestComponentWorkspace)
	assert.Equal(t, 9, len(tenant1Ue2DevTestTestComponentDeps))
	assert.Equal(t, "catalog/terraform/services/service-1", tenant1Ue2DevTestTestComponentDeps[0])
	assert.Equal(t, "catalog/terraform/services/service-2", tenant1Ue2DevTestTestComponentDeps[1])
	assert.Equal(t, "catalog/terraform/spacelift-and-backend-override-1", tenant1Ue2DevTestTestComponentDeps[2])
	assert.Equal(t, "catalog/terraform/test-component", tenant1Ue2DevTestTestComponentDeps[3])
	assert.Equal(t, "mixins/region/us-east-2", tenant1Ue2DevTestTestComponentDeps[4])
	assert.Equal(t, "mixins/stage/dev", tenant1Ue2DevTestTestComponentDeps[5])
	assert.Equal(t, "orgs/cp/_defaults", tenant1Ue2DevTestTestComponentDeps[6])
	assert.Equal(t, "orgs/cp/tenant1/_defaults", tenant1Ue2DevTestTestComponentDeps[7])
	assert.Equal(t, "orgs/cp/tenant1/dev/us-east-2", tenant1Ue2DevTestTestComponentDeps[8])
	assert.Equal(t, "tenant1-ue2-dev", tenant1Ue2DevTestTestComponentTerraformWorkspace)

	var tenant1Ue2DevTestTestComponent2 map[string]any
	component = "test/test-component"
	tenant := "tenant1"
	environment := "ue2"
	stage := "dev"
	tenant1Ue2DevTestTestComponent2, err = c.ProcessComponentFromContext(&c.ComponentFromContextParams{
		Component:   component,
		Namespace:   "",
		Tenant:      tenant,
		Environment: environment,
		Stage:       stage,
	})
	assert.Nil(t, err)
	tenant1Ue2DevTestTestComponentBackend2 := tenant1Ue2DevTestTestComponent2["backend"].(map[string]any)
	tenant1Ue2DevTestTestComponentRemoteStateBackend2 := tenant1Ue2DevTestTestComponent2["remote_state_backend"].(map[string]any)
	tenant1Ue2DevTestTestComponentBaseComponent2 := tenant1Ue2DevTestTestComponent2["component"]
	tenant1Ue2DevTestTestComponentTerraformWorkspace2 := tenant1Ue2DevTestTestComponent2["workspace"].(string)
	tenant1Ue2DevTestTestComponentWorkspace2 := tenant1Ue2DevTestTestComponent2["workspace"].(string)
	tenant1Ue2DevTestTestComponentBackendWorkspaceKeyPrefix2 := tenant1Ue2DevTestTestComponentBackend2["workspace_key_prefix"].(string)
	tenant1Ue2DevTestTestComponentRemoteStateBackendWorkspaceKeyPrefix2 := tenant1Ue2DevTestTestComponentRemoteStateBackend2["workspace_key_prefix"].(string)
	tenant1Ue2DevTestTestComponentDeps2 := tenant1Ue2DevTestTestComponent2["deps"].([]any)
	assert.Equal(t, "test-test-component", tenant1Ue2DevTestTestComponentBackendWorkspaceKeyPrefix2)
	assert.Equal(t, "test-test-component", tenant1Ue2DevTestTestComponentRemoteStateBackendWorkspaceKeyPrefix2)
	assert.Equal(t, "test/test-component", tenant1Ue2DevTestTestComponentBaseComponent2)
	assert.Equal(t, "tenant1-ue2-dev", tenant1Ue2DevTestTestComponentWorkspace2)
	assert.Equal(t, 9, len(tenant1Ue2DevTestTestComponentDeps2))
	assert.Equal(t, "catalog/terraform/services/service-1", tenant1Ue2DevTestTestComponentDeps2[0])
	assert.Equal(t, "catalog/terraform/services/service-2", tenant1Ue2DevTestTestComponentDeps2[1])
	assert.Equal(t, "catalog/terraform/spacelift-and-backend-override-1", tenant1Ue2DevTestTestComponentDeps2[2])
	assert.Equal(t, "catalog/terraform/test-component", tenant1Ue2DevTestTestComponentDeps2[3])
	assert.Equal(t, "mixins/region/us-east-2", tenant1Ue2DevTestTestComponentDeps2[4])
	assert.Equal(t, "mixins/stage/dev", tenant1Ue2DevTestTestComponentDeps2[5])
	assert.Equal(t, "orgs/cp/_defaults", tenant1Ue2DevTestTestComponentDeps2[6])
	assert.Equal(t, "orgs/cp/tenant1/_defaults", tenant1Ue2DevTestTestComponentDeps2[7])
	assert.Equal(t, "orgs/cp/tenant1/dev/us-east-2", tenant1Ue2DevTestTestComponentDeps2[8])
	assert.Equal(t, "tenant1-ue2-dev", tenant1Ue2DevTestTestComponentTerraformWorkspace2)

	yamlConfig, err = u.ConvertToYAML(tenant1Ue2DevTestTestComponent)
	assert.Nil(t, err)
	t.Log(yamlConfig)

	var tenant1Ue2DevTestTestComponentOverrideComponent map[string]any
	component = "test/test-component-override"
	stack = "tenant1-ue2-dev"
	tenant1Ue2DevTestTestComponentOverrideComponent, err = c.ProcessComponentInStack(component, stack, "", "")
	assert.Nil(t, err)
	tenant1Ue2DevTestTestComponentOverrideComponentBackend := tenant1Ue2DevTestTestComponentOverrideComponent["backend"].(map[string]any)
	tenant1Ue2DevTestTestComponentOverrideComponentBaseComponent := tenant1Ue2DevTestTestComponentOverrideComponent["component"].(string)
	tenant1Ue2DevTestTestComponentOverrideComponentWorkspace := tenant1Ue2DevTestTestComponentOverrideComponent["workspace"].(string)
	tenant1Ue2DevTestTestComponentOverrideComponentBackendWorkspaceKeyPrefix := tenant1Ue2DevTestTestComponentOverrideComponentBackend["workspace_key_prefix"].(string)
	tenant1Ue2DevTestTestComponentOverrideComponentDeps := tenant1Ue2DevTestTestComponentOverrideComponent["deps"].([]any)
	tenant1Ue2DevTestTestComponentOverrideComponentRemoteStateBackend := tenant1Ue2DevTestTestComponentOverrideComponent["remote_state_backend"].(map[string]any)
	tenant1Ue2DevTestTestComponentOverrideComponentRemoteStateBackendVal2 := tenant1Ue2DevTestTestComponentOverrideComponentRemoteStateBackend["val2"].(string)
	assert.Equal(t, "test-test-component", tenant1Ue2DevTestTestComponentOverrideComponentBackendWorkspaceKeyPrefix)
	assert.Equal(t, "test/test-component", tenant1Ue2DevTestTestComponentOverrideComponentBaseComponent)
	assert.Equal(t, "test-component-override-workspace-override", tenant1Ue2DevTestTestComponentOverrideComponentWorkspace)

	assert.Equal(t, 10, len(tenant1Ue2DevTestTestComponentOverrideComponentDeps))
	assert.Equal(t, "catalog/terraform/services/service-1-override", tenant1Ue2DevTestTestComponentOverrideComponentDeps[0])
	assert.Equal(t, "catalog/terraform/services/service-2-override", tenant1Ue2DevTestTestComponentOverrideComponentDeps[1])
	assert.Equal(t, "catalog/terraform/spacelift-and-backend-override-1", tenant1Ue2DevTestTestComponentOverrideComponentDeps[2])
	assert.Equal(t, "catalog/terraform/test-component", tenant1Ue2DevTestTestComponentOverrideComponentDeps[3])
	assert.Equal(t, "catalog/terraform/test-component-override", tenant1Ue2DevTestTestComponentOverrideComponentDeps[4])
	assert.Equal(t, "mixins/region/us-east-2", tenant1Ue2DevTestTestComponentOverrideComponentDeps[5])
	assert.Equal(t, "mixins/stage/dev", tenant1Ue2DevTestTestComponentOverrideComponentDeps[6])
	assert.Equal(t, "orgs/cp/_defaults", tenant1Ue2DevTestTestComponentOverrideComponentDeps[7])
	assert.Equal(t, "orgs/cp/tenant1/_defaults", tenant1Ue2DevTestTestComponentOverrideComponentDeps[8])
	assert.Equal(t, "orgs/cp/tenant1/dev/us-east-2", tenant1Ue2DevTestTestComponentOverrideComponentDeps[9])

	assert.Equal(t, "2", tenant1Ue2DevTestTestComponentOverrideComponentRemoteStateBackendVal2)

	var tenant1Ue2DevTestTestComponentOverrideComponent2 map[string]any
	component = "test/test-component-override-2"
	tenant = "tenant1"
	environment = "ue2"
	stage = "dev"
	tenant1Ue2DevTestTestComponentOverrideComponent2, err = c.ProcessComponentFromContext(&c.ComponentFromContextParams{
		Component:   component,
		Namespace:   "",
		Tenant:      tenant,
		Environment: environment,
		Stage:       stage,
	})
	assert.Nil(t, err)
	tenant1Ue2DevTestTestComponentOverrideComponent2Backend := tenant1Ue2DevTestTestComponentOverrideComponent2["backend"].(map[string]any)
	tenant1Ue2DevTestTestComponentOverrideComponent2Workspace := tenant1Ue2DevTestTestComponentOverrideComponent2["workspace"].(string)
	tenant1Ue2DevTestTestComponentOverrideComponent2WorkspaceKeyPrefix := tenant1Ue2DevTestTestComponentOverrideComponent2Backend["workspace_key_prefix"].(string)
	assert.Equal(t, "tenant1-ue2-dev-test-test-component-override-2", tenant1Ue2DevTestTestComponentOverrideComponent2Workspace)
	assert.Equal(t, "test-test-component", tenant1Ue2DevTestTestComponentOverrideComponent2WorkspaceKeyPrefix)

	yamlConfig, err = u.ConvertToYAML(tenant1Ue2DevTestTestComponentOverrideComponent2)
	assert.Nil(t, err)
	t.Log(yamlConfig)

	// Test having a dash `-` in the stage name
	var tenant1Ue2Test1TestTestComponentOverrideComponent2 map[string]any
	component = "test/test-component-override-2"
	tenant = "tenant1"
	environment = "ue2"
	stage = "test-1"
	tenant1Ue2Test1TestTestComponentOverrideComponent2, err = c.ProcessComponentFromContext(&c.ComponentFromContextParams{
		Component:   component,
		Namespace:   "",
		Tenant:      tenant,
		Environment: environment,
		Stage:       stage,
	})
	assert.Nil(t, err)
	tenant1Ue2Test1TestTestComponentOverrideComponent2Backend := tenant1Ue2Test1TestTestComponentOverrideComponent2["backend"].(map[string]any)
	tenant1Ue2Test1TestTestComponentOverrideComponent2Workspace := tenant1Ue2Test1TestTestComponentOverrideComponent2["workspace"].(string)
	tenant1Ue2Test1TestTestComponentOverrideComponent2WorkspaceKeyPrefix := tenant1Ue2Test1TestTestComponentOverrideComponent2Backend["workspace_key_prefix"].(string)
	assert.Equal(t, "tenant1-ue2-test-1-test-test-component-override-2", tenant1Ue2Test1TestTestComponentOverrideComponent2Workspace)
	assert.Equal(t, "test-test-component", tenant1Ue2Test1TestTestComponentOverrideComponent2WorkspaceKeyPrefix)

	var tenant1Ue2DevTestTestComponentOverrideComponent3 map[string]any
	component = "test/test-component-override-3"
	stack = "tenant1-ue2-dev"
	tenant1Ue2DevTestTestComponentOverrideComponent3, err = c.ProcessComponentInStack(component, stack, "", "")
	assert.Nil(t, err)

	tenant1Ue2DevTestTestComponentOverrideComponent3Deps := tenant1Ue2DevTestTestComponentOverrideComponent3["deps"].([]any)

	assert.Equal(t, 11, len(tenant1Ue2DevTestTestComponentOverrideComponent3Deps))
	assert.Equal(t, "catalog/terraform/mixins/test-2", tenant1Ue2DevTestTestComponentOverrideComponent3Deps[0])
	assert.Equal(t, "catalog/terraform/services/service-1-override-2", tenant1Ue2DevTestTestComponentOverrideComponent3Deps[1])
	assert.Equal(t, "catalog/terraform/services/service-2-override-2", tenant1Ue2DevTestTestComponentOverrideComponent3Deps[2])
	assert.Equal(t, "catalog/terraform/spacelift-and-backend-override-1", tenant1Ue2DevTestTestComponentOverrideComponent3Deps[3])
	assert.Equal(t, "catalog/terraform/test-component", tenant1Ue2DevTestTestComponentOverrideComponent3Deps[4])
	assert.Equal(t, "catalog/terraform/test-component-override-3", tenant1Ue2DevTestTestComponentOverrideComponent3Deps[5])
	assert.Equal(t, "mixins/region/us-east-2", tenant1Ue2DevTestTestComponentOverrideComponent3Deps[6])
	assert.Equal(t, "mixins/stage/dev", tenant1Ue2DevTestTestComponentOverrideComponent3Deps[7])
	assert.Equal(t, "orgs/cp/_defaults", tenant1Ue2DevTestTestComponentOverrideComponent3Deps[8])
	assert.Equal(t, "orgs/cp/tenant1/_defaults", tenant1Ue2DevTestTestComponentOverrideComponent3Deps[9])
	assert.Equal(t, "orgs/cp/tenant1/dev/us-east-2", tenant1Ue2DevTestTestComponentOverrideComponent3Deps[10])
}

// TestComponentProcessorConsistency verifies that ProcessComponentInStack and
// ProcessComponentFromContext return identical results for the same component and stack.
// This is critical because the provider uses both paths depending on user input.
func TestComponentProcessorConsistency(t *testing.T) {
	component := "test/test-component"
	stack := "tenant1-ue2-dev"

	resultByStack, err := c.ProcessComponentInStack(component, stack, "", "")
	require.NoError(t, err)
	require.NotNil(t, resultByStack)

	resultByContext, err := c.ProcessComponentFromContext(&c.ComponentFromContextParams{
		Component:   component,
		Namespace:   "",
		Tenant:      "tenant1",
		Environment: "ue2",
		Stage:       "dev",
	})
	require.NoError(t, err)
	require.NotNil(t, resultByContext)

	
	// Both paths should produce the same backend config
	stackBackend := resultByStack["backend"].(map[string]any)
	contextBackend := resultByContext["backend"].(map[string]any)
	assert.Equal(t, stackBackend["workspace_key_prefix"], contextBackend["workspace_key_prefix"])

	// Both paths should produce the same workspace
	assert.Equal(t, resultByStack["workspace"], resultByContext["workspace"])

	// Both paths should produce the same base component
	assert.Equal(t, resultByStack["component"], resultByContext["component"])

	// Both paths should produce the same deps
	assert.Equal(t, resultByStack["deps"], resultByContext["deps"])
}

// TestComponentProcessorProdStack tests processing a component in a different stack (prod)
// to ensure the provider works across environments.
func TestComponentProcessorProdStack(t *testing.T) {
	component := "top-level-component1"
	stack := "tenant1-ue2-prod"

	result, err := c.ProcessComponentInStack(component, stack, "", "")
	require.NoError(t, err)
	require.NotNil(t, result)

	workspace := result["workspace"].(string)
	assert.Equal(t, "tenant1-ue2-prod", workspace)

	backend := result["backend"].(map[string]any)
	assert.Equal(t, "top-level-component1", backend["workspace_key_prefix"])

	vars := result["vars"].(map[string]any)
	assert.Equal(t, "tenant1", vars["tenant"])
	assert.Equal(t, "ue2", vars["environment"])
	assert.Equal(t, "prod", vars["stage"])
}

// TestComponentProcessorFromContextProdStack tests ProcessComponentFromContext
// for the prod stack to verify tenant/environment/stage resolution.
func TestComponentProcessorFromContextProdStack(t *testing.T) {
	result, err := c.ProcessComponentFromContext(&c.ComponentFromContextParams{
		Component:   "top-level-component1",
		Namespace:   "",
		Tenant:      "tenant1",
		Environment: "ue2",
		Stage:       "prod",
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	workspace := result["workspace"].(string)
	assert.Equal(t, "tenant1-ue2-prod", workspace)

	vars := result["vars"].(map[string]any)
	assert.Equal(t, "tenant1", vars["tenant"])
	assert.Equal(t, "ue2", vars["environment"])
	assert.Equal(t, "prod", vars["stage"])
}

// TestComponentProcessorFromContextNilParams tests that ProcessComponentFromContext
// returns an error when called with nil params.
func TestComponentProcessorFromContextNilParams(t *testing.T) {
	_, err := c.ProcessComponentFromContext(nil)
	assert.NotNil(t, err)
}

// TestComponentProcessorInfraVpc tests the infra/vpc component which is commonly
// used with remote-state modules (the original bug report scenario).
func TestComponentProcessorInfraVpc(t *testing.T) {
	component := "infra/vpc"
	stack := "tenant1-ue2-dev"

	result, err := c.ProcessComponentInStack(component, stack, "", "")
	require.NoError(t, err)
	require.NotNil(t, result)

	backend := result["backend"].(map[string]any)
	assert.Equal(t, "infra-vpc", backend["workspace_key_prefix"])
	assert.Equal(t, "s3", result["backend_type"])
}

// TestComponentProcessorWithProcessingDisabled verifies that passing
// WithProcessTemplates(false) and WithProcessYamlFunctions(false) still returns
// valid component configuration (backend, workspace, vars). This is the mode
// used by the provider to avoid spawning child processes.
func TestComponentProcessorWithProcessingDisabled(t *testing.T) {
	component := "test/test-component"
	stack := "tenant1-ue2-dev"

	result, err := c.ProcessComponentInStack(component, stack, "", "",
		c.WithProcessTemplates(false),
		c.WithProcessYamlFunctions(false),
	)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Backend config should still be present
	backend := result["backend"].(map[string]any)
	assert.Equal(t, "test-test-component", backend["workspace_key_prefix"])

	// Workspace should still be present
	assert.Equal(t, "tenant1-ue2-dev", result["workspace"])

	// Vars should still be present
	vars := result["vars"].(map[string]any)
	assert.Equal(t, "tenant1", vars["tenant"])
	assert.Equal(t, "ue2", vars["environment"])
	assert.Equal(t, "dev", vars["stage"])

	// Base component should still be present
	assert.Equal(t, "test/test-component", result["component"])
}

// TestComponentProcessorFromContextWithProcessingDisabled verifies that
// ProcessComponentFromContext also works with processing disabled.
func TestComponentProcessorFromContextWithProcessingDisabled(t *testing.T) {
	result, err := c.ProcessComponentFromContext(&c.ComponentFromContextParams{
		Component:   "test/test-component",
		Namespace:   "",
		Tenant:      "tenant1",
		Environment: "ue2",
		Stage:       "dev",
	},
		c.WithProcessTemplates(false),
		c.WithProcessYamlFunctions(false),
	)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Backend config should still be present
	backend := result["backend"].(map[string]any)
	assert.Equal(t, "test-test-component", backend["workspace_key_prefix"])

	// Workspace should still be present
	assert.Equal(t, "tenant1-ue2-dev", result["workspace"])

	// Vars should still be present
	vars := result["vars"].(map[string]any)
	assert.Equal(t, "tenant1", vars["tenant"])
	assert.Equal(t, "ue2", vars["environment"])
	assert.Equal(t, "dev", vars["stage"])
}

// TestComponentProcessorDisabledMatchesEnabled verifies that for components
// without templates or YAML functions, the results are identical regardless
// of whether processing is enabled or disabled.
func TestComponentProcessorDisabledMatchesEnabled(t *testing.T) {
	component := "test/test-component"
	stack := "tenant1-ue2-dev"

	resultEnabled, err := c.ProcessComponentInStack(component, stack, "", "")
	require.NoError(t, err)

	resultDisabled, err := c.ProcessComponentInStack(component, stack, "", "",
		c.WithProcessTemplates(false),
		c.WithProcessYamlFunctions(false),
	)
	require.NoError(t, err)

	// For components without templates or YAML functions, key fields should match
	assert.Equal(t, resultEnabled["workspace"], resultDisabled["workspace"])
	assert.Equal(t, resultEnabled["component"], resultDisabled["component"])
	assert.Equal(t, resultEnabled["backend_type"], resultDisabled["backend_type"])

	enabledBackend := resultEnabled["backend"].(map[string]any)
	disabledBackend := resultDisabled["backend"].(map[string]any)
	assert.Equal(t, enabledBackend["workspace_key_prefix"], disabledBackend["workspace_key_prefix"])

	enabledVars := resultEnabled["vars"].(map[string]any)
	disabledVars := resultDisabled["vars"].(map[string]any)
	assert.Equal(t, enabledVars["tenant"], disabledVars["tenant"])
	assert.Equal(t, enabledVars["environment"], disabledVars["environment"])
	assert.Equal(t, enabledVars["stage"], disabledVars["stage"])
}

// TestComponentProcessorConsistencyWithProcessingDisabled verifies that both
// ProcessComponentInStack and ProcessComponentFromContext return the same
// results when processing is disabled — the mode the provider actually uses.
func TestComponentProcessorConsistencyWithProcessingDisabled(t *testing.T) {
	component := "infra/vpc"
	stack := "tenant1-ue2-dev"

	resultByStack, err := c.ProcessComponentInStack(component, stack, "", "",
		c.WithProcessTemplates(false),
		c.WithProcessYamlFunctions(false),
	)
	require.NoError(t, err)

	resultByContext, err := c.ProcessComponentFromContext(&c.ComponentFromContextParams{
		Component:   component,
		Namespace:   "",
		Tenant:      "tenant1",
		Environment: "ue2",
		Stage:       "dev",
	},
		c.WithProcessTemplates(false),
		c.WithProcessYamlFunctions(false),
	)
	require.NoError(t, err)

	// Both paths should produce the same results with processing disabled
	assert.Equal(t, resultByStack["workspace"], resultByContext["workspace"])
	assert.Equal(t, resultByStack["component"], resultByContext["component"])
	assert.Equal(t, resultByStack["backend_type"], resultByContext["backend_type"])

	stackBackend := resultByStack["backend"].(map[string]any)
	contextBackend := resultByContext["backend"].(map[string]any)
	assert.Equal(t, stackBackend["workspace_key_prefix"], contextBackend["workspace_key_prefix"])
}

// TestComponentProcessorWithEmptyBasePath verifies that passing empty string
// for atmosBasePath (the default provider behavior) works correctly after the
// base path resolution changes in Atmos v1.210.1. Empty base path should
// trigger git root -> config dir -> CWD fallback chain.
func TestComponentProcessorWithEmptyBasePath(t *testing.T) {
	component := "test/test-component"
	stack := "tenant1-ue2-dev"

	// Empty atmosBasePath is the default provider behavior
	result, err := c.ProcessComponentInStack(component, stack, "", "")
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify all key fields are present and correct
	backend := result["backend"].(map[string]any)
	assert.Equal(t, "test-test-component", backend["workspace_key_prefix"])
	assert.Equal(t, "tenant1-ue2-dev", result["workspace"])
	assert.Equal(t, "test/test-component", result["component"])

	vars := result["vars"].(map[string]any)
	assert.Equal(t, "tenant1", vars["tenant"])
	assert.Equal(t, "ue2", vars["environment"])
	assert.Equal(t, "dev", vars["stage"])
}

// TestComponentProcessorFromContextWithEmptyBasePath verifies that
// ProcessComponentFromContext with empty AtmosBasePath works correctly
// after the base path resolution changes in Atmos v1.210.1.
func TestComponentProcessorFromContextWithEmptyBasePath(t *testing.T) {
	result, err := c.ProcessComponentFromContext(&c.ComponentFromContextParams{
		Component:          "test/test-component",
		Namespace:          "",
		Tenant:             "tenant1",
		Environment:        "ue2",
		Stage:              "dev",
		AtmosCliConfigPath: "",
		AtmosBasePath:      "",
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	backend := result["backend"].(map[string]any)
	assert.Equal(t, "test-test-component", backend["workspace_key_prefix"])
	assert.Equal(t, "tenant1-ue2-dev", result["workspace"])

	vars := result["vars"].(map[string]any)
	assert.Equal(t, "tenant1", vars["tenant"])
	assert.Equal(t, "ue2", vars["environment"])
	assert.Equal(t, "dev", vars["stage"])
}

// TestComponentProcessorWithRelativeBasePath verifies that the test fixture's
// relative base_path ("../../examples/tests") in atmos.yaml still resolves
// correctly after the path resolution changes. This validates that config-file
// relative paths (non-runtime source) continue to resolve relative to the
// atmos.yaml directory.
func TestComponentProcessorWithRelativeBasePath(t *testing.T) {
	// The atmos.yaml in internal/component/ has base_path: "../../examples/tests"
	// This is a config-file relative path and should resolve relative to atmos.yaml location
	component := "infra/vpc"
	stack := "tenant1-ue2-dev"

	result, err := c.ProcessComponentInStack(component, stack, "", "")
	require.NoError(t, err)
	require.NotNil(t, result)

	backend := result["backend"].(map[string]any)
	assert.Equal(t, "infra-vpc", backend["workspace_key_prefix"])
	assert.Equal(t, "s3", result["backend_type"])
	assert.Equal(t, "tenant1-ue2-dev", result["workspace"])
}

// TestComponentProcessorFromContextWithRuntimeRelativeBasePath verifies that
// passing a relative path via AtmosBasePath (runtime source) resolves correctly
// relative to CWD, not the config directory. This covers the core fix path from
// Atmos v1.210.1.
//
// To prove CWD-based resolution, the test changes CWD to the repo root and passes
// "examples/tests" as AtmosBasePath with an explicit AtmosCliConfigPath pointing to
// the original config directory. If AtmosBasePath resolved relative to the config
// directory (internal/component/), the path "internal/component/examples/tests" would
// not exist and the test would fail. It only passes because runtime-source paths
// resolve relative to CWD (repo root), where "examples/tests" does exist.
func TestComponentProcessorFromContextWithRuntimeRelativeBasePath(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	repoRoot := filepath.Clean(filepath.Join(cwd, "../.."))
	t.Chdir(repoRoot)

	result, err := c.ProcessComponentFromContext(&c.ComponentFromContextParams{
		Component:          "infra/vpc",
		Namespace:          "",
		Tenant:             "tenant1",
		Environment:        "ue2",
		Stage:              "dev",
		AtmosCliConfigPath: filepath.Join(repoRoot, "internal", "component"),
		AtmosBasePath:      "examples/tests", // runtime relative — resolves from CWD (repo root), not config dir
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	backend := result["backend"].(map[string]any)
	assert.Equal(t, "infra-vpc", backend["workspace_key_prefix"])
	assert.Equal(t, "tenant1-ue2-dev", result["workspace"])
}

// TestComponentProcessorWithProcessingDisabledAndEmptyPath verifies the actual
// provider mode: both template/YAML function processing disabled AND empty base
// path. This is the exact combination used in production by the provider's
// dataSourceComponentConfigRead function.
func TestComponentProcessorWithProcessingDisabledAndEmptyPath(t *testing.T) {
	component := "test/test-component"
	stack := "tenant1-ue2-dev"

	result, err := c.ProcessComponentInStack(component, stack, "", "",
		c.WithProcessTemplates(false),
		c.WithProcessYamlFunctions(false),
	)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Backend config must be present for remote-state lookup
	backend := result["backend"].(map[string]any)
	assert.Equal(t, "test-test-component", backend["workspace_key_prefix"])

	// Workspace must be present for remote-state lookup
	assert.Equal(t, "tenant1-ue2-dev", result["workspace"])

	// Vars must be present
	vars := result["vars"].(map[string]any)
	assert.Equal(t, "tenant1", vars["tenant"])
	assert.Equal(t, "ue2", vars["environment"])
	assert.Equal(t, "dev", vars["stage"])

	// Base component must be present
	assert.Equal(t, "test/test-component", result["component"])

	// Also test ProcessComponentFromContext in the same mode
	resultCtx, err := c.ProcessComponentFromContext(&c.ComponentFromContextParams{
		Component:          component,
		Namespace:          "",
		Tenant:             "tenant1",
		Environment:        "ue2",
		Stage:              "dev",
		AtmosCliConfigPath: "",
		AtmosBasePath:      "",
	},
		c.WithProcessTemplates(false),
		c.WithProcessYamlFunctions(false),
	)
	require.NoError(t, err)
	require.NotNil(t, resultCtx)

	// Both paths should produce the same results
	assert.Equal(t, result["workspace"], resultCtx["workspace"])
	assert.Equal(t, result["component"], resultCtx["component"])

	ctxBackend := resultCtx["backend"].(map[string]any)
	assert.Equal(t, backend["workspace_key_prefix"], ctxBackend["workspace_key_prefix"])
}

func TestComponentProcessorHierarchicalInheritance(t *testing.T) {
	var yamlConfig string
	component := "derived-component-2"
	tenant := "tenant1"
	environment := "ue2"
	stage := "test-1"

	componentMap, err := c.ProcessComponentFromContext(&c.ComponentFromContextParams{
		Component:   component,
		Namespace:   "",
		Tenant:      tenant,
		Environment: environment,
		Stage:       stage,
	})
	assert.Nil(t, err)

	componentVars := componentMap["vars"].(map[string]any)
	componentHierarchicalInheritanceTestVar := componentVars["hierarchical_inheritance_test"].(string)
	assert.Equal(t, "base-component-1", componentHierarchicalInheritanceTestVar)

	yamlConfig, err = u.ConvertToYAML(componentMap)
	assert.Nil(t, err)
	t.Log(yamlConfig)
}

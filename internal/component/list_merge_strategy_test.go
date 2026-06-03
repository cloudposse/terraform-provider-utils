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

// TestComponentProcessorListMergeStrategy verifies that a component-level
// `settings.list_merge_strategy` is honored when Atmos deep-merges a component's
// `vars` across inheritance layers (Atmos #2480, released in v1.220.0).
//
// All three components below override the same base list with the same override
// list; only their `settings.list_merge_strategy` differs. The strategy therefore
// fully determines the merged result:
//
//	replace (global default) -> override replaces base
//	append                   -> base elements followed by override elements
//	merge                    -> element-wise deep merge of list items
//
// The provider exposes this path through utils_component_config /
// utils_describe_stacks (ProcessComponentInStack), so this test pins the
// component-level override behavior that the embedded Atmos library provides.
func TestComponentProcessorListMergeStrategy(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Use the isolated fixture so this test cannot perturb (or be perturbed by)
	// the shared examples/tests fixture used by the other component tests.
	fixture := filepath.Join(cwd, "testdata", "list-merge-strategy")
	t.Chdir(fixture)

	const stack = "acme-ue2-test"

	process := []c.ProcessOption{
		c.WithProcessTemplates(false),
		c.WithProcessYamlFunctions(false),
	}

	t.Run("replace overrides the base list", func(t *testing.T) {
		result, err := c.ProcessComponentInStack("list-replace", stack, "", "", process...)
		require.NoError(t, err)

		vars := result["vars"].(map[string]any)
		assert.Equal(t, []any{"x", "y", "z"}, vars["my_list"])
	})

	t.Run("append concatenates base then override", func(t *testing.T) {
		result, err := c.ProcessComponentInStack("list-append", stack, "", "", process...)
		require.NoError(t, err)

		vars := result["vars"].(map[string]any)
		assert.Equal(t, []any{"a", "b", "x", "y", "z"}, vars["my_list"])
	})

	t.Run("merge deep-merges list items element-wise", func(t *testing.T) {
		result, err := c.ProcessComponentInStack("list-merge", stack, "", "", process...)
		require.NoError(t, err)

		vars := result["vars"].(map[string]any)
		list := vars["my_map_list"].([]any)
		require.Len(t, list, 1, "element-wise merge keeps a single merged item")

		item := list[0].(map[string]any)
		assert.Equal(t, "v1b", item["key1"], "overridden key wins")
		assert.Equal(t, "v2", item["key2"], "base-only key is preserved")
		assert.Equal(t, "v3", item["key3"], "override-only key is added")

		yamlConfig, err := u.ConvertToYAML(result)
		require.NoError(t, err)
		t.Log(yamlConfig)
	})
}

package merge

import (
	"testing"

	"github.com/stretchr/testify/assert"

	m "github.com/cloudposse/atmos/pkg/merge"
	"github.com/cloudposse/atmos/pkg/schema"
	u "github.com/cloudposse/atmos/pkg/utils"
)

func TestMergeBasic(t *testing.T) {
	cliConfig := schema.AtmosConfiguration{}

	map1 := map[string]any{"foo": "bar"}
	map2 := map[string]any{"baz": "bat"}

	inputs := []map[string]any{map1, map2}
	expected := map[string]any{"foo": "bar", "baz": "bat"}

	result, err := m.Merge(&cliConfig, inputs)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestMergeBasicOverride(t *testing.T) {
	cliConfig := schema.AtmosConfiguration{}

	map1 := map[string]any{"foo": "bar"}
	map2 := map[string]any{"baz": "bat"}
	map3 := map[string]any{"foo": "ood"}

	inputs := []map[string]any{map1, map2, map3}
	expected := map[string]any{"foo": "ood", "baz": "bat"}

	result, err := m.Merge(&cliConfig, inputs)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestMergeListReplace(t *testing.T) {
	cliConfig := schema.AtmosConfiguration{
		Settings: schema.AtmosSettings{
			ListMergeStrategy: m.ListMergeStrategyReplace,
		},
	}

	map1 := map[string]any{
		"list": []string{"1", "2", "3"},
	}

	map2 := map[string]any{
		"list": []string{"4", "5", "6"},
	}

	inputs := []map[string]any{map1, map2}
	expected := map[string]any{"list": []any{"4", "5", "6"}}

	result, err := m.Merge(&cliConfig, inputs)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)

	yamlConfig, err := u.ConvertToYAML(result)
	assert.Nil(t, err)
	t.Log(yamlConfig)
}

func TestMergeListAppend(t *testing.T) {
	cliConfig := schema.AtmosConfiguration{
		Settings: schema.AtmosSettings{
			ListMergeStrategy: m.ListMergeStrategyAppend,
		},
	}

	map1 := map[string]any{
		"list": []string{"1", "2", "3"},
	}

	map2 := map[string]any{
		"list": []string{"4", "5", "6"},
	}

	inputs := []map[string]any{map1, map2}
	expected := map[string]any{"list": []any{"1", "2", "3", "4", "5", "6"}}

	result, err := m.Merge(&cliConfig, inputs)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)

	yamlConfig, err := u.ConvertToYAML(result)
	assert.Nil(t, err)
	t.Log(yamlConfig)
}

func TestMergeListMerge(t *testing.T) {
	cliConfig := schema.AtmosConfiguration{
		Settings: schema.AtmosSettings{
			ListMergeStrategy: m.ListMergeStrategyMerge,
		},
	}

	map1 := map[string]any{
		"list": []map[string]string{
			{
				"1": "1",
				"2": "2",
				"3": "3",
				"4": "4",
			},
		},
	}

	map2 := map[string]any{
		"list": []map[string]string{
			{
				"1": "1b",
				"2": "2",
				"3": "3b",
				"5": "5",
			},
		},
	}

	inputs := []map[string]any{map1, map2}

	result, err := m.Merge(&cliConfig, inputs)
	assert.Nil(t, err)

	var mergedList []any
	var ok bool

	if mergedList, ok = result["list"].([]any); !ok {
		t.Errorf("invalid merge result: %v", result)
	}

	merged := mergedList[0].(map[string]any)

	assert.Equal(t, "1b", merged["1"])
	assert.Equal(t, "2", merged["2"])
	assert.Equal(t, "3b", merged["3"])
	assert.Equal(t, "4", merged["4"])
	assert.Equal(t, "5", merged["5"])

	yamlConfig, err := u.ConvertToYAML(result)
	assert.Nil(t, err)
	t.Log(yamlConfig)
}

// TestMergeWithOptionsNilConfig tests MergeWithOptions with nil AtmosConfiguration,
// which is how the provider's deep_merge_json and deep_merge_yaml data sources call it.
func TestMergeWithOptionsNilConfig(t *testing.T) {
	map1 := map[string]any{"foo": "bar"}
	map2 := map[string]any{"baz": "bat"}

	inputs := []map[string]any{map1, map2}
	expected := map[string]any{"foo": "bar", "baz": "bat"}

	result, err := m.MergeWithOptions(nil, inputs, false, false)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

// TestMergeWithOptionsNilConfigOverride tests that MergeWithOptions correctly overrides
// values when called with nil AtmosConfiguration.
func TestMergeWithOptionsNilConfigOverride(t *testing.T) {
	map1 := map[string]any{"key": "original", "keep": "this"}
	map2 := map[string]any{"key": "override", "new": "value"}

	inputs := []map[string]any{map1, map2}
	expected := map[string]any{"key": "override", "keep": "this", "new": "value"}

	result, err := m.MergeWithOptions(nil, inputs, false, false)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

// TestMergeWithOptionsNilConfigNestedMaps tests deep merge with nested maps
// using nil AtmosConfiguration.
func TestMergeWithOptionsNilConfigNestedMaps(t *testing.T) {
	map1 := map[string]any{
		"top": map[string]any{
			"nested1": "value1",
			"nested2": "value2",
		},
	}
	map2 := map[string]any{
		"top": map[string]any{
			"nested2": "override",
			"nested3": "value3",
		},
	}

	inputs := []map[string]any{map1, map2}

	result, err := m.MergeWithOptions(nil, inputs, false, false)
	assert.Nil(t, err)

	top := result["top"].(map[string]any)
	assert.Equal(t, "value1", top["nested1"])
	assert.Equal(t, "override", top["nested2"])
	assert.Equal(t, "value3", top["nested3"])
}

// TestMergeWithOptionsAppendList tests MergeWithOptions with appendSlice=true
// and nil AtmosConfiguration.
func TestMergeWithOptionsAppendList(t *testing.T) {
	map1 := map[string]any{
		"list": []any{"a", "b"},
	}
	map2 := map[string]any{
		"list": []any{"c", "d"},
	}

	inputs := []map[string]any{map1, map2}

	result, err := m.MergeWithOptions(nil, inputs, true, false)
	assert.Nil(t, err)

	list := result["list"].([]any)
	assert.Equal(t, 4, len(list))
	assert.Equal(t, "a", list[0])
	assert.Equal(t, "b", list[1])
	assert.Equal(t, "c", list[2])
	assert.Equal(t, "d", list[3])
}

// TestMergeWithOptionsDeepCopyList tests MergeWithOptions with deepCopyList=true
// and nil AtmosConfiguration, exercising element-wise list merging.
func TestMergeWithOptionsDeepCopyList(t *testing.T) {
	map1 := map[string]any{
		"items": []any{
			map[string]any{"name": "a", "value": "1"},
		},
	}
	map2 := map[string]any{
		"items": []any{
			map[string]any{"name": "b", "value": "2"},
		},
	}

	inputs := []map[string]any{map1, map2}

	result, err := m.MergeWithOptions(nil, inputs, false, true)
	assert.Nil(t, err)

	items, ok := result["items"].([]any)
	assert.True(t, ok, "items should be a slice")
	// deepCopyList=true merges element-by-element (element[0] with element[0]),
	// so one merged element should result from two single-element lists
	assert.Equal(t, 1, len(items), "expected one merged item from element-wise merge")

	merged := items[0].(map[string]any)
	assert.Equal(t, "b", merged["name"], "name should be overridden by second input")
	assert.Equal(t, "2", merged["value"], "value should be overridden by second input")
}

// TestMergeWithOptionsSingleInput tests MergeWithOptions with a single input map.
func TestMergeWithOptionsSingleInput(t *testing.T) {
	map1 := map[string]any{"key": "value"}

	inputs := []map[string]any{map1}

	result, err := m.MergeWithOptions(nil, inputs, false, false)
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"key": "value"}, result)
}

// TestMergeWithOptionsEmptyInputs tests MergeWithOptions with empty inputs.
func TestMergeWithOptionsEmptyInputs(t *testing.T) {
	inputs := []map[string]any{}

	result, err := m.MergeWithOptions(nil, inputs, false, false)
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{}, result)
}

// TestMergeWithOptionsTypeOverride verifies the deep-merge "src always wins at the
// leaf level regardless of type" contract through the public API that the provider's
// utils_deep_merge_json and utils_deep_merge_yaml data sources call
// (m.MergeWithOptions(nil, inputs, false, false)).
//
// Atmos v1.212.0 reworked the deep-merge engine (mergo -> reflection-free native merge,
// atmos #2201) and fixed a regression where guards rejected some heterogeneous-type
// overrides (atmos #2248). Those overrides — a list replaced by a map, a scalar, a bool
// or null — show up in real-world configs (e.g. blanking out an inherited list with `{}`
// or null). This test pins that behavior so a future atmos bump that reintroduces the
// guard is caught here rather than silently changing user merge results.
func TestMergeWithOptionsTypeOverride(t *testing.T) {
	tests := []struct {
		name    string
		first   map[string]any
		second  map[string]any
		wantKey string
		wantVal any
	}{
		{
			name:    "list overridden by empty map",
			first:   map[string]any{"accounts": []any{"a", "b"}},
			second:  map[string]any{"accounts": map[string]any{}},
			wantKey: "accounts",
			wantVal: map[string]any{},
		},
		{
			name:    "list overridden by non-empty map",
			first:   map[string]any{"accounts": []any{"a", "b"}},
			second:  map[string]any{"accounts": map[string]any{"only": "one"}},
			wantKey: "accounts",
			wantVal: map[string]any{"only": "one"},
		},
		{
			name:    "list overridden by scalar string",
			first:   map[string]any{"cidrs": []any{"10.0.0.0/16"}},
			second:  map[string]any{"cidrs": "10.99.0.0/16"},
			wantKey: "cidrs",
			wantVal: "10.99.0.0/16",
		},
		{
			name:    "list overridden by scalar int",
			first:   map[string]any{"ports": []any{80, 443}},
			second:  map[string]any{"ports": 8080},
			wantKey: "ports",
			wantVal: 8080,
		},
		{
			name:    "list overridden by null",
			first:   map[string]any{"rules": []any{"allow"}},
			second:  map[string]any{"rules": nil},
			wantKey: "rules",
			wantVal: nil,
		},
		{
			name:    "list overridden by bool false",
			first:   map[string]any{"features": []any{"a", "b"}},
			second:  map[string]any{"features": false},
			wantKey: "features",
			wantVal: false,
		},
		{
			name:    "map overridden by scalar string",
			first:   map[string]any{"backend": map[string]any{"type": "s3"}},
			second:  map[string]any{"backend": "local"},
			wantKey: "backend",
			wantVal: "local",
		},
		{
			name:    "scalar overridden by list",
			first:   map[string]any{"cidr": "10.0.0.0/16"},
			second:  map[string]any{"cidr": []any{"10.0.0.0/16", "10.1.0.0/16"}},
			wantKey: "cidr",
			wantVal: []any{"10.0.0.0/16", "10.1.0.0/16"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputs := []map[string]any{tt.first, tt.second}

			// Matches the data sources' call: nil config, no append, no deep-copy list.
			result, err := m.MergeWithOptions(nil, inputs, false, false)
			assert.Nil(t, err, "type override must not error")
			assert.Equal(t, tt.wantVal, result[tt.wantKey])
		})
	}
}

// TestMergeWithOptionsTypeOverrideNestedWithSliceFlags verifies that a heterogeneous-type
// override of a nested value holds even when the slice flags (appendSlice / sliceDeepCopy)
// are set — i.e. the type override takes precedence over list-merge handling, matching the
// utils_deep_merge_* data sources when callers pass append=true or deep_copy_list=true.
func TestMergeWithOptionsTypeOverrideNestedWithSliceFlags(t *testing.T) {
	flagCases := []struct {
		name          string
		appendSlice   bool
		sliceDeepCopy bool
	}{
		{"appendSlice", true, false},
		{"sliceDeepCopy", false, true},
		{"bothFlags", true, true},
	}

	for _, fc := range flagCases {
		t.Run(fc.name, func(t *testing.T) {
			first := map[string]any{
				"vars": map[string]any{"accounts": []any{"a", "b"}},
			}
			second := map[string]any{
				"vars": map[string]any{"accounts": map[string]any{}},
			}

			result, err := m.MergeWithOptions(nil, []map[string]any{first, second}, fc.appendSlice, fc.sliceDeepCopy)
			assert.Nil(t, err, "type override must work regardless of slice flags")

			vars := result["vars"].(map[string]any)
			assert.Equal(t, map[string]any{}, vars["accounts"])
		})
	}
}

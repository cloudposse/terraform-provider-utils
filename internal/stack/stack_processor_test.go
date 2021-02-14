package stack

import (
	c "github.com/cloudposse/terraform-provider-utils/internal/convert"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestStackProcessor(t *testing.T) {
	filePaths := []string{
		"../../examples/data-sources/utils_stack_config_yaml/stacks/uw2-dev.yaml",
	}

	yamlResult, err := ProcessYAMLConfigFiles(filePaths)
	assert.Nil(t, err)
	assert.Equal(t, len(yamlResult), 1)

	mapResult, err := c.YAMLToMapOfInterfaces(yamlResult[0])
	assert.Nil(t, err)

	terraformComponents := mapResult["components"].(map[interface{}]interface{})["terraform"].(map[interface{}]interface{})
	helmfileComponents := mapResult["components"].(map[interface{}]interface{})["helmfile"].(map[interface{}]interface{})

	auroraPostgres2Component := terraformComponents["aurora-postgres-2"].(map[interface{}]interface{})
	assert.Equal(t, auroraPostgres2Component["component"], "aurora-postgres")
	assert.Equal(t, auroraPostgres2Component["settings"].(map[interface{}]interface{})["spacelift"].(map[interface{}]interface{})["branch"], "dev")
	assert.Equal(t, auroraPostgres2Component["vars"].(map[interface{}]interface{})["instance_type"], "db.r4.xlarge")
	assert.Equal(t, auroraPostgres2Component["env"].(map[interface{}]interface{})["ENV_TEST_1"].(string), "test1_override2")
	assert.Equal(t, auroraPostgres2Component["env"].(map[interface{}]interface{})["ENV_TEST_2"].(string), "test2_override2")
	assert.Equal(t, auroraPostgres2Component["env"].(map[interface{}]interface{})["ENV_TEST_3"].(string), "test3")
	assert.Equal(t, auroraPostgres2Component["env"].(map[interface{}]interface{})["ENV_TEST_4"].(string), "test4")
	assert.Equal(t, auroraPostgres2Component["env"].(map[interface{}]interface{})["ENV_TEST_5"].(string), "test5")
	assert.Equal(t, auroraPostgres2Component["env"].(map[interface{}]interface{})["ENV_TEST_6"].(string), "test6")
	assert.Equal(t, auroraPostgres2Component["env"].(map[interface{}]interface{})["ENV_TEST_7"].(string), "test7")
	assert.Equal(t, auroraPostgres2Component["env"].(map[interface{}]interface{})["ENV_TEST_8"].(string), "test8")
	assert.Nil(t, auroraPostgres2Component["env"].(map[interface{}]interface{})["ENV_TEST_9"])

	eksComponent := terraformComponents["eks"].(map[interface{}]interface{})
	assert.Equal(t, eksComponent["settings"].(map[interface{}]interface{})["spacelift"].(map[interface{}]interface{})["workspace_enabled"], true)
	assert.Equal(t, eksComponent["settings"].(map[interface{}]interface{})["spacelift"].(map[interface{}]interface{})["branch"], "test")
	assert.Equal(t, eksComponent["vars"].(map[interface{}]interface{})["spotinst_oceans"].(map[interface{}]interface{})["main"].(map[interface{}]interface{})["max_group_size"], 3)
	assert.Equal(t, eksComponent["vars"].(map[interface{}]interface{})["spotinst_instance_profile"], "eg-gbl-dev-spotinst-worker")
	assert.Equal(t, eksComponent["env"].(map[interface{}]interface{})["ENV_TEST_1"].(string), "test1_override")
	assert.Equal(t, eksComponent["env"].(map[interface{}]interface{})["ENV_TEST_2"].(string), "test2_override")
	assert.Equal(t, eksComponent["env"].(map[interface{}]interface{})["ENV_TEST_3"].(string), "test3")
	assert.Equal(t, eksComponent["env"].(map[interface{}]interface{})["ENV_TEST_4"].(string), "test4")
	assert.Nil(t, eksComponent["env"].(map[interface{}]interface{})["ENV_TEST_5"])

	accountComponent := terraformComponents["account"].(map[interface{}]interface{})
	assert.Equal(t, accountComponent["backend_type"].(string), "s3")
	assert.Equal(t, accountComponent["backend"].(map[interface{}]interface{})["workspace_key_prefix"], "account")
	assert.Equal(t, accountComponent["backend"].(map[interface{}]interface{})["bucket"], "eg-uw2-root-tfstate")
	assert.Nil(t, accountComponent["backend"].(map[interface{}]interface{})["role_arn"])

	datadogHelmfileComponent := helmfileComponents["datadog"].(map[interface{}]interface{})
	assert.Equal(t, datadogHelmfileComponent["vars"].(map[interface{}]interface{})["account_number"], "1234567890")
	assert.Equal(t, datadogHelmfileComponent["vars"].(map[interface{}]interface{})["installed"], true)
	assert.Equal(t, datadogHelmfileComponent["vars"].(map[interface{}]interface{})["stage"], "dev")
	assert.Equal(t, datadogHelmfileComponent["vars"].(map[interface{}]interface{})["processAgent"].(map[interface{}]interface{})["enabled"], true)

	yamlConfig, err := yaml.Marshal(mapResult)
	assert.Nil(t, err)
	t.Log(string(yamlConfig))
}

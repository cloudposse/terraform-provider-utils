package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"

	m "github.com/cloudposse/terraform-provider-utils/internal/merge"
)

func dataSourceDeepMergeYaml() *schema.Resource {
	return &schema.Resource{
		Description: "The `deep_merge_yaml` data source accepts a list of YAML strings as input and deep merges into a single YAML string as output.",

		ReadContext: dataSourceDeepMergeYamlRead,

		Schema: map[string]*schema.Schema{
			"inputs": {
				Description: "A list of arbitrary maps that is deep merged into the `output` attribute.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
			},
			"output": {
				Description: "The deep-merged map.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceDeepMergeYamlRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inputs := d.Get("inputs")

	inputMaps := make([]map[string]interface{}, 0)
	for _, current := range inputs.([]interface{}) {
		var data map[string]interface{}
		byt := []byte(current.(string))

		if err := yaml.Unmarshal(byt, &data); err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[DEBUG] current data: %v", data)
		inputMaps = append(inputMaps, data)
	}

	result, err := m.Merge(inputMaps)
	log.Printf("[DEBUG] merged data: %v", result)
	if err != nil {
		return diag.FromErr(err)
	}

	yamlResult, err := yaml.Marshal(result)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] ***result: %v", string(yamlResult))

	d.Set("output", string(yamlResult))

	d.SetId("static")

	return nil
}

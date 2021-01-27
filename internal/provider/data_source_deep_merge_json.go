package provider

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	m "github.com/cloudposse/terraform-provider-utils/internal/merge"
)

func dataSourceDeepMergeJson() *schema.Resource {
	return &schema.Resource{
		Description: "The `deep_merge_json` data source accepts a list of JSON strings as input and deep merges into a single JSON string as output.",

		ReadContext: dataSourceDeepMergeJsonRead,

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

func dataSourceDeepMergeJsonRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inputs := d.Get("inputs")

	inputMaps := make([]map[string]interface{}, 0)
	for _, current := range inputs.([]interface{}) {
		var data map[string]interface{}
		byt := []byte(current.(string))

		if err := json.Unmarshal(byt, &data); err != nil {
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

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] ***result: %v", string(jsonResult))

	d.Set("output", string(jsonResult))

	d.SetId("static")

	return nil
}

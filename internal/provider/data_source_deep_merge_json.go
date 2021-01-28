package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/cloudposse/terraform-provider-utils/internal/convert"
)

func dataSourceDeepMergeJSON() *schema.Resource {
	return &schema.Resource{
		Description: "The `deep_merge_json` data source accepts a list of JSON strings as input and deep merges into a single JSON string as output.",

		ReadContext: dataSourceDeepMergeJSONRead,

		Schema: map[string]*schema.Schema{
			"inputs": {
				Description: "A list JSON strings that is deep merged into the `output` attribute.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
			},
			"output": {
				Description: "The deep-merged output.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceDeepMergeJSONRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inputs := d.Get("inputs")

	data, err := c.JSONSliceOfInterfaceToSliceOfMaps(inputs.([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	// Convert result to JSON
	jsonResult, err := json.Marshal(data)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("output", string(jsonResult))
	d.SetId("static")

	return nil
}

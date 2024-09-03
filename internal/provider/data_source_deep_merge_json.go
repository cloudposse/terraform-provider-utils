package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	jsoniter "github.com/json-iterator/go"

	c "github.com/cloudposse/atmos/pkg/convert"
	m "github.com/cloudposse/atmos/pkg/merge"
)

func dataSourceDeepMergeJSON() *schema.Resource {
	return &schema.Resource{
		Description: "The `deep_merge_json` data source accepts a list of JSON strings as input and deep merges into a single JSON string as output.",

		ReadContext: dataSourceDeepMergeJSONRead,

		Schema: map[string]*schema.Schema{
			"input": {
				Description: "A list of JSON strings that is deep merged into the `output` attribute.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
			},
			"append_list": {
				Description: "A boolean flag to enable/disable appending lists instead of overwriting them.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"deep_copy_list": {
				Description: "A boolean flag to enable/disable merging of list elements one by one.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"output": {
				Description: "The deep-merged output.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceDeepMergeJSONRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	input := d.Get("input")
	appendList := d.Get("append_list").(bool)
	deepCopyList := d.Get("deep_copy_list").(bool)

	data, err := JSONSliceOfInterfaceToSliceOfMaps(input.([]any))
	if err != nil {
		return diag.FromErr(err)
	}

	merged, err := m.MergeWithOptions(data, appendList, deepCopyList)
	if err != nil {
		return diag.FromErr(err)
	}

	// Convert result to JSON
	var json = jsoniter.ConfigDefault
	jsonResult, err := json.Marshal(merged)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("output", string(jsonResult))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(c.MakeId(jsonResult))

	return nil
}

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDeepMerge() *schema.Resource {
	return &schema.Resource{
		Description: "The `deep_merge` data source accepts a list of maps as input and deep merges them as output.",

		ReadContext: dataSourceDeepMergeRead,

		Schema: map[string]*schema.Schema{
			"inputs": {
				Description: "A listx of arbitrary maps that is deep merged into the `output` attribute.",
				Type:        schema.TypeMap,
				Required:    true,
			},
			"output": {
				Description: "The deep-merged map.",
				Type:        schema.TypeMap,
				Computed:    true,
			},
		},
	}
}

func dataSourceDeepMergeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inputs := d.Get("inputs")
	d.Set("output", inputs)

	d.SetId("static")

	// diag.Errorf("an error")
	return nil
}

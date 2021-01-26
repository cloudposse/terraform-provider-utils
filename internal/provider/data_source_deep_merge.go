package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	m "github.com/cloudposse/terraform-provider-utils/internal/merge"
)

func dataSourceDeepMerge() *schema.Resource {
	return &schema.Resource{
		Description: "The `deep_merge` data source accepts a list of maps as input and deep merges them as output.",

		ReadContext: dataSourceDeepMergeRead,

		Schema: map[string]*schema.Schema{
			"inputs": {
				Description: "A list of arbitrary maps that is deep merged into the `output` attribute.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeMap},
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

	result, err := m.Merge(inputs.([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("output", result)

	d.SetId("static")

	return nil
}

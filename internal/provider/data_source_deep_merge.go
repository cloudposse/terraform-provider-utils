package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDeepMerge() *schema.Resource {
	return &schema.Resource{
		Description: "The `deep_merge` data source accepts a list of maps as input and deep merges them as output.",

		Read: dataSourceDeepMergeRead,

		Schema: map[string]*schema.Schema{
			"inputs": {
				Description: "A listx of arbitrary maps that is deep merged into the `output` attribute.",
				Type:        schema.TypeMap,
				Optional:    true,
			},
			"output": {
				Description: "The deep-merged map.",
				Type:        schema.TypeMap,
				Computed:    true,
			},
		},
	}
}

func dataSourceDeepMergeRead(d *schema.ResourceData, meta interface{}) error {

	inputs := d.Get("inputs")
	d.Set("output", inputs)

	d.SetId("static")
	return nil
}

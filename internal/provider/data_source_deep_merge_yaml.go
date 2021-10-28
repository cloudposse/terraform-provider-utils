package provider

import (
	"context"
	c "github.com/cloudposse/terraform-provider-utils/internal/convert"
	m "github.com/cloudposse/terraform-provider-utils/internal/merge"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"
)

func dataSourceDeepMergeYAML() *schema.Resource {
	return &schema.Resource{
		Description: "The `deep_merge_yaml` data source accepts a list of YAML strings as input and deep merges into a single YAML string as output.",

		ReadContext: dataSourceDeepMergeYAMLRead,

		Schema: map[string]*schema.Schema{
			"input": {
				Description: "A list of YAML strings that is deep merged into the `output` attribute.",
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

func dataSourceDeepMergeYAMLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	input := d.Get("input")
	appendList := d.Get("append_list").(bool)
	deepCopyList := d.Get("deep_copy_list").(bool)

	data, err := c.YAMLSliceOfInterfaceToSliceOfMaps(input.([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	merged, err := m.MergeWithOptions(data, appendList, deepCopyList)
	if err != nil {
		return diag.FromErr(err)
	}

	// Convert result to YAML
	yamlResult, err := yaml.Marshal(merged)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("output", string(yamlResult))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(c.MakeId(yamlResult))

	return nil
}

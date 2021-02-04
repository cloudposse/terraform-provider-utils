package provider

import (
	"context"

	c "github.com/cloudposse/terraform-provider-utils/internal/convert"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"
)

func dataSourceStackConfigYAML() *schema.Resource {
	return &schema.Resource{
		Description: "The `stack_config_yaml` data source accepts a list of file names " +
			"as input and returns a single YAML string with stack configurations as output.",

		ReadContext: dataSourceStackConfigYAMLRead,

		Schema: map[string]*schema.Schema{
			"inputs": {
				Description: "A list file names that is processed into the `output` attribute.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
			},
			"output": {
				Description: "The stack config output.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceStackConfigYAMLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inputs := d.Get("inputs")

	data, err := c.YAMLSliceOfInterfaceToSliceOfMaps(inputs.([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	// Convert result to YAML
	yamlResult, err := yaml.Marshal(data)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("output", string(yamlResult))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("static")

	return nil
}

package provider

import (
	"context"

	s "github.com/cloudposse/terraform-provider-utils/internal/stack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"
)

func dataSourceStackConfigYAML() *schema.Resource {
	return &schema.Resource{
		Description: "The `stack_config_yaml` data source accepts a list of stack config file names " +
			"and returns a list of stack configurations.",

		ReadContext: dataSourceStackConfigYAMLRead,

		Schema: map[string]*schema.Schema{
			"input": {
				Description: "A list of stack config file names.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
			},
			"output": {
				Description: "A list of stack configurations.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
		},
	}
}

func dataSourceStackConfigYAMLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := d.Get("input")
	paths := input.([]interface{})
	result := make([]string, len(paths))

	for _, path := range paths {
		config, err := s.ProcessYAMLConfigFile(path.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		yamlConfig, err := yaml.Marshal(config)
		if err != nil {
			return diag.FromErr(err)
		}
		result = append(result, string(yamlConfig))
	}

	err := d.Set("output", result)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("static")

	return nil
}

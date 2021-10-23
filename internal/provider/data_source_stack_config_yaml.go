package provider

import (
	"context"
	s "github.com/cloudposse/atmos/pkg/stack"
	c "github.com/cloudposse/terraform-provider-utils/internal/convert"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
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
			"base_path": {
				Description: "Stack config base path.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"process_stack_deps": {
				Description: "A boolean flag to enable/disable processing all stack dependencies for the components.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"process_component_deps": {
				Description: "A boolean flag to enable/disable processing config dependencies for the components.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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
	processStackDeps := d.Get("process_stack_deps")
	processComponentDeps := d.Get("process_component_deps")
	basePath := d.Get("base_path")

	paths, err := c.SliceOfInterfacesToSliceOfStrings(input.([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	result, _, err := s.ProcessYAMLConfigFiles(
		basePath.(string),
		paths,
		processStackDeps.(bool),
		processComponentDeps.(bool))

	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("output", result)
	if err != nil {
		return diag.FromErr(err)
	}

	id := c.MakeId([]byte(strings.Join(result, "")))
	d.SetId(id)

	return nil
}

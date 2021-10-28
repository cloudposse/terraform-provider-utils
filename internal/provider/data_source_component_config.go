package provider

import (
	"context"
	c "github.com/cloudposse/terraform-provider-utils/internal/convert"
	s "github.com/cloudposse/terraform-provider-utils/internal/stack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func dataSourceComponentConfig() *schema.Resource {
	return &schema.Resource{
		Description: "The `component_config` data source accepts a component and a stack name " +
			"and returns the component configuration in the stack",

		ReadContext: dataSourceComponentConfigRead,

		Schema: map[string]*schema.Schema{
			"component": {
				Description: "Component name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"stack": {
				Description: "Stack name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"stack_name_pattern": {
				Description: "Stack name pattern.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"base_path": {
				Description: "Stack config base path.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"process_stack_deps": {
				Description: "A boolean flag to enable/disable processing all stack dependencies for the component.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"process_component_deps": {
				Description: "A boolean flag to enable/disable processing config dependencies for the component.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"output": {
				Description: "Component configuration.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceComponentConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

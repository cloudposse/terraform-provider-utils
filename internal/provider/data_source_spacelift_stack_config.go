package provider

import (
	"context"
	c "github.com/cloudposse/atmos/pkg/convert"
	s "github.com/cloudposse/atmos/pkg/spacelift"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"
)

func dataSourceSpaceliftStackConfig() *schema.Resource {
	return &schema.Resource{
		Description: "The `spacelift_stack_config` data source accepts a list of stack config file names " +
			"and returns a map of Spacelift stack configurations.",

		ReadContext: dataSourceSpaceliftStackConfigRead,

		Schema: map[string]*schema.Schema{
			"input": {
				Description: "A list of stack config file names.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
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
			"process_imports": {
				Description: "A boolean flag to enable/disable processing stack imports.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"stack_config_path_template": {
				Description: "Stack config path template.",
				Type:        schema.TypeString,
				Required:    true,
			},
			// https://www.terraform.io/plugin/sdkv2/schemas/schema-types#typemap
			"env": {
				Description: "Map of ENV vars in the format 'key=value'. These ENV vars will be set before executing the data source",
				Type:        schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Default:  nil,
			},
			"output": {
				Description: "A map of Spacelift stack configurations.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceSpaceliftStackConfigRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	input := d.Get("input")
	processStackDeps := d.Get("process_stack_deps")
	processComponentDeps := d.Get("process_component_deps")
	processImports := d.Get("process_imports")
	stackConfigPathTemplate := d.Get("stack_config_path_template")
	stacksBasePath := d.Get("base_path")
	env := d.Get("env").(map[string]any)

	err := setEnv(env)
	if err != nil {
		return diag.FromErr(err)
	}

	paths, err := c.SliceOfInterfacesToSliceOfStrings(input.([]any))
	if err != nil {
		return diag.FromErr(err)
	}

	spaceliftStacks, err := s.CreateSpaceliftStacks(
		stacksBasePath.(string),
		"",
		"",
		paths,
		processStackDeps.(bool),
		processComponentDeps.(bool),
		processImports.(bool),
		stackConfigPathTemplate.(string))

	if err != nil {
		return diag.FromErr(err)
	}

	yamlConfig, err := yaml.Marshal(spaceliftStacks)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("output", string(yamlConfig))
	if err != nil {
		return diag.FromErr(err)
	}

	id := c.MakeId(yamlConfig)
	d.SetId(id)

	return nil
}

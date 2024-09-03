package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	cfg "github.com/cloudposse/atmos/pkg/config"
	c "github.com/cloudposse/atmos/pkg/convert"
	atmosSchema "github.com/cloudposse/atmos/pkg/schema"
	s "github.com/cloudposse/atmos/pkg/stack"
	u "github.com/cloudposse/atmos/pkg/utils"
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
				Description: "A list of stack configurations.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
		},
	}
}

func dataSourceStackConfigYAMLRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	cliConfig, err := cfg.InitCliConfig(atmosSchema.ConfigAndStacksInfo{}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	input := d.Get("input")
	processStackDeps := d.Get("process_stack_deps")
	processComponentDeps := d.Get("process_component_deps")
	stacksBasePath := d.Get("base_path")
	env := d.Get("env").(map[string]any)

	err = setEnv(env)
	if err != nil {
		return diag.FromErr(err)
	}

	paths, err := u.SliceOfInterfacesToSliceOfStrings(input.([]any))
	if err != nil {
		return diag.FromErr(err)
	}

	result, _, _, err := s.ProcessYAMLConfigFiles(
		cliConfig,
		stacksBasePath.(string),
		"",
		"",
		paths,
		processStackDeps.(bool),
		processComponentDeps.(bool),
		false,
	)

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

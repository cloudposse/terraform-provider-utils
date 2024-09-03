package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v3"

	cfg "github.com/cloudposse/atmos/pkg/config"
	c "github.com/cloudposse/atmos/pkg/convert"
	"github.com/cloudposse/atmos/pkg/describe"
	s "github.com/cloudposse/atmos/pkg/schema"
)

func dataSourceDescribeStacks() *schema.Resource {
	return &schema.Resource{
		Description: "The `describe_stacks` data source shows configuration for Atmos stacks and components in the stacks",

		ReadContext: dataSourceDescribeStacksRead,

		Schema: map[string]*schema.Schema{
			"stack": {
				Description: "Atmos stack to filter by.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"tenant": {
				Description: "Tenant to filter by.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"namespace": {
				Description: "Namespace to filter by.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"environment": {
				Description: "Environment to filter by.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"stage": {
				Description: "Stage to filter by.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"components": {
				Description: "List of Atmos components to filter by.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
			},
			"component_types": {
				Description: "List of component types to filter by.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
			},
			"sections": {
				Description: "Output only the specified component sections.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
			},
			"ignore_errors": {
				Description: "Flag to ignore errors in the provider when executing 'describe stacks' command.",
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
			"atmos_cli_config_path": {
				Description: "Atmos CLI config path.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"atmos_base_path": {
				Description: "Atmos base path to components and stacks.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"output": {
				Description: "Stack configurations.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceDescribeStacksRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	stack := d.Get("stack").(string)
	namespace := d.Get("namespace").(string)
	tenant := d.Get("tenant").(string)
	environment := d.Get("environment").(string)
	stage := d.Get("stage").(string)
	components := d.Get("components")
	componentTypes := d.Get("component_types")
	sections := d.Get("sections")
	ignoreErrors := d.Get("ignore_errors").(bool)
	env := d.Get("env").(map[string]any)
	atmosCliConfigPath := d.Get("atmos_cli_config_path").(string)
	atmosBasePath := d.Get("atmos_base_path").(string)

	var result map[string]any
	var err error
	var yamlConfig []byte

	componentsList, err := SliceOfInterfacesToSliceOfStrings(components.([]any))
	if err != nil {
		return diag.FromErr(err)
	}

	componentTypesList, err := SliceOfInterfacesToSliceOfStrings(componentTypes.([]any))
	if err != nil {
		return diag.FromErr(err)
	}

	sectionsList, err := SliceOfInterfacesToSliceOfStrings(sections.([]any))
	if err != nil {
		return diag.FromErr(err)
	}

	err = setEnv(env)
	if err != nil {
		return diag.FromErr(err)
	}

	info := s.ConfigAndStacksInfo{
		AtmosBasePath:      atmosBasePath,
		AtmosCliConfigPath: atmosCliConfigPath,
	}

	cliConfig, err := cfg.InitCliConfig(info, true)
	if err != nil {
		return diag.FromErr(err)
	}

	var filterByStack string

	if stack != "" {
		filterByStack = stack
	} else if namespace != "" || tenant != "" || environment != "" || stage != "" {
		filterByStack, err = cfg.GetStackNameFromContextAndStackNamePattern(namespace, tenant, environment, stage, cliConfig.Stacks.NamePattern)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	result, err = describe.ExecuteDescribeStacks(
		cliConfig,
		filterByStack,
		componentsList,
		componentTypesList,
		sectionsList,
		false)
	if err != nil && !ignoreErrors {
		return diag.FromErr(err)
	}

	if err != nil {
		result = map[string]any{}
	}

	yamlConfig, err = yaml.Marshal(result)
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

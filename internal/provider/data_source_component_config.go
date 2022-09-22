package provider

import (
	"context"
	p "github.com/cloudposse/atmos/pkg/component"
	c "github.com/cloudposse/atmos/pkg/convert"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"
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
				Optional:    true,
				Default:     "",
			},
			"tenant": {
				Description: "Tenant.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"namespace": {
				Description: "Namespace.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"environment": {
				Description: "Environment.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"stage": {
				Description: "Stage.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"ignore_errors": {
				Description: "Flag to ignore errors if the component is not found in the stack.",
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
				Description: "Component configuration.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceComponentConfigRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	component := d.Get("component").(string)
	stack := d.Get("stack").(string)
	namespace := d.Get("namespace").(string)
	tenant := d.Get("tenant").(string)
	environment := d.Get("environment").(string)
	stage := d.Get("stage").(string)
	ignoreErrors := d.Get("ignore_errors").(bool)
	env := d.Get("env").(map[string]any)

	var result map[string]any
	var err error
	var yamlConfig []byte

	err = setEnv(env)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(stack) > 0 {
		result, err = p.ProcessComponentInStack(component, stack)
		if err != nil && !ignoreErrors {
			return diag.FromErr(err)
		}
	} else {
		result, err = p.ProcessComponentFromContext(component, namespace, tenant, environment, stage)
		if err != nil && !ignoreErrors {
			return diag.FromErr(err)
		}
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

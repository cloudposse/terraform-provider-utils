package provider

import (
	"context"
	p "github.com/cloudposse/terraform-provider-utils/internal/component"
	c "github.com/cloudposse/terraform-provider-utils/internal/convert"
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
			"output": {
				Description: "Component configuration.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceComponentConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	component := d.Get("component").(string)
	stack := d.Get("stack").(string)
	tenant := d.Get("tenant").(string)
	environment := d.Get("environment").(string)
	stage := d.Get("stage").(string)
	var result map[string]interface{}
	var err error
	var yamlConfig []byte

	if len(stack) > 0 {
		result, err = p.ProcessComponentInStack(component, stack)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		result, err = p.ProcessComponentFromContext(component, tenant, environment, stage)
		if err != nil {
			return diag.FromErr(err)
		}
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
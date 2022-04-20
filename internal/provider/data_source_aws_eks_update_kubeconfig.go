package provider

import (
	"context"
	p "github.com/cloudposse/atmos/pkg/component"
	c "github.com/cloudposse/atmos/pkg/convert"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"
)

func dataSourceAwsEksUpdateKubeconfig() *schema.Resource {
	return &schema.Resource{
		Description: "The `component_config` data source accepts a component and a stack name " +
			"and returns the component configuration in the stack",

		ReadContext: dataSourceAwsEksUpdateKubeconfigRead,

		// https://docs.aws.amazon.com/cli/latest/reference/eks/update-kubeconfig.html
		Schema: map[string]*schema.Schema{
			"component": {
				Description: "Component name.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"stack": {
				Description: "Stack name.",
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
			"profile": {
				Description: "AWS profile to use for cluster authentication.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"cluster_name": {
				Description: "EKS cluster name.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"kubeconfig": {
				Description: "kubeconfig file path to write the kubeconfig to. By default, the configuration is written to the first file path in the KUBECONFIG environment variable (if it is set) or the default kubeconfig path (.kube/config) in your home directory",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"role-arn": {
				Description: "IAM role to assume for cluster authentication.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"alias": {
				Description: "Alias for the cluster context name. Defaults to match cluster ARN.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
		},
	}
}

func dataSourceAwsEksUpdateKubeconfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if err != nil {
		result = map[string]interface{}{}
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

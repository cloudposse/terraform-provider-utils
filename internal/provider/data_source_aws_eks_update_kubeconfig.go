package provider

import (
	"context"
	a "github.com/cloudposse/atmos/pkg/aws"
	g "github.com/cloudposse/atmos/pkg/config"
	c "github.com/cloudposse/atmos/pkg/convert"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func dataSourceAwsEksUpdateKubeconfig() *schema.Resource {
	return &schema.Resource{
		Description: "The 'utils_aws_eks_update_kubeconfig' data source executes 'aws eks update-kubeconfig' commands",

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
			"role_arn": {
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
			"region": {
				Description: "AWS region.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"output": {
				Description: "Output.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAwsEksUpdateKubeconfigRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	component := d.Get("component").(string)
	stack := d.Get("stack").(string)
	tenant := d.Get("tenant").(string)
	environment := d.Get("environment").(string)
	stage := d.Get("stage").(string)
	profile := d.Get("profile").(string)
	clusterName := d.Get("cluster_name").(string)
	kubeconfig := d.Get("kubeconfig").(string)
	roleArn := d.Get("role_arn").(string)
	alias := d.Get("alias").(string)
	region := d.Get("region").(string)

	kubeconfigContext := g.AwsEksUpdateKubeconfigContext{
		Component:   component,
		Stack:       stack,
		Profile:     profile,
		ClusterName: clusterName,
		Region:      region,
		Kubeconfig:  kubeconfig,
		RoleArn:     roleArn,
		Alias:       alias,
		Tenant:      tenant,
		Environment: environment,
		Stage:       stage,
	}

	err := a.ExecuteAwsEksUpdateKubeconfig(kubeconfigContext)
	if err != nil {
		return diag.FromErr(err)
	}

	utc := time.Now().UTC().String()

	err = d.Set("output", utc)
	if err != nil {
		return diag.FromErr(err)
	}

	id := c.MakeId([]byte(utc))
	d.SetId(id)

	return nil
}

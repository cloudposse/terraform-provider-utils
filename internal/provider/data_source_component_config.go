package provider

import (
	"context"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v3"

	c "github.com/cloudposse/atmos/pkg/convert"
	p "github.com/cloudposse/atmos/pkg/describe"
)

// parseBoolEnv reads an environment variable and parses it as a boolean.
// Returns defaultVal if the variable is unset or has an unrecognized value.
func parseBoolEnv(envVar string, defaultVal bool) bool {
	val := os.Getenv(envVar)
	if val == "" {
		return defaultVal
	}
	switch strings.ToLower(val) {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	default:
		return defaultVal
	}
}

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
			"atmos_cli_config_path": {
				Description: "atmos CLI config path.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"atmos_base_path": {
				Description: "atmos base path to components and stacks.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"process_templates": {
				Description: "Enable Go template processing in the component config output. " +
					"Defaults to true. Can also be set via ATMOS_PROCESS_TEMPLATES env var.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"process_yaml_functions": {
				Description: "Enable YAML function processing (e.g., !terraform.output) in the component config output. " +
					"Defaults to false to avoid ETXTBSY crashes from child process execution inside the provider. " +
					"Can also be set via ATMOS_PROCESS_FUNCTIONS env var.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
	atmosCliConfigPath := d.Get("atmos_cli_config_path").(string)
	atmosBasePath := d.Get("atmos_base_path").(string)
	// Default from env var, can be overridden by schema attribute
	processTemplates := parseBoolEnv("ATMOS_PROCESS_TEMPLATES", true)
	if v, ok := d.GetOk("process_templates"); ok {
		processTemplates = v.(bool)
	}

	processYamlFunctions := parseBoolEnv("ATMOS_PROCESS_FUNCTIONS", false)
	if v, ok := d.GetOk("process_yaml_functions"); ok {
		processYamlFunctions = v.(bool)
	}

	var result map[string]any
	var err error
	var yamlConfig []byte

	if env != nil {
		err = setEnv(env)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	atmosMu.Lock()
	if len(stack) > 0 {
		result, err = p.ProcessComponentInStack(component, stack, atmosCliConfigPath, atmosBasePath,
			p.WithProcessTemplates(processTemplates),
			p.WithProcessYamlFunctions(processYamlFunctions),
		)
	} else {
		result, err = p.ProcessComponentFromContext(&p.ComponentFromContextParams{
			Component:          component,
			Namespace:          namespace,
			Tenant:             tenant,
			Environment:        environment,
			Stage:              stage,
			AtmosCliConfigPath: atmosCliConfigPath,
			AtmosBasePath:      atmosBasePath,
		},
			p.WithProcessTemplates(processTemplates),
			p.WithProcessYamlFunctions(processYamlFunctions),
		)
	}
	atmosMu.Unlock()

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

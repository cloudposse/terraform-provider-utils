package provider

import (
	"math/rand"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
	rand.Seed(time.Now().Unix())
}

// New returns a *schema.Provider.
func New() *schema.Provider {
	return &schema.Provider{
		Schema:       map[string]*schema.Schema{},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"cloudposse_deep_merge": dataSourceDeepMerge(),
		},
	}
}

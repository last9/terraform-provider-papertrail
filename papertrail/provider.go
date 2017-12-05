package papertrail

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:         schema.TypeString,
				Description:  "X-Papertrail-Token",
				InputDefault: "",
				Sensitive:    true,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("PAPERTRAIL_TOKEN", ""),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"papertrail_system":       resourcePapertrailSystem(),
			"papertrail_group":        resourcePapertrailGroup(),
			"papertrail_system_group": resourcePapertrailSystemGroup(),
			"papertrail_search":       resourcePapertrailSearch(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"papertrail_user":            dataSourcePapertrailUser(),
			"papertrail_log_destination": dataSourcePapertrailLogDestination(),
		},
		ConfigureFunc: providerConfigure,
	}
}

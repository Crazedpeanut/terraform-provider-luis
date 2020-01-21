package luis

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/crazedpeanut/terraform-provider-luis/luis/internal"
	"github.com/crazedpeanut/terraform-provider-luis/luis/internal/clients"
	"github.com/hashicorp/terraform/terraform"
)

// Provider of luis resources
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ConfigureFunc: configureFunc,
		Schema: map[string]*schema.Schema{
			"authoring_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				Default:  "westus.api.cognitive.microsoft.com",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"luis_application": internal.ResourceApplication(),
			"luis_version": internal.ResourceVersion(),
		},
	}
}

func configureFunc(d *schema.ResourceData) (interface{}, error) {
	options := clients.ClientOptions{
		AuthoringKey: d.Get("authoring_key").(string),
		Domain:       d.Get("domain").(string),
	}

	return clients.NewClient(&options), nil
}

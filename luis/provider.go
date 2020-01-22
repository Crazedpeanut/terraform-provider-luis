package luis

import (
	"github.com/crazedpeanut/terraform-provider-luis/luis/internal"
	"github.com/crazedpeanut/terraform-provider-luis/luis/internal/clients"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider of luis resources
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ConfigureFunc: configureFunc,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "westus.api.cognitive.microsoft.com",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"luis_application": internal.ResourceApplication(),
			"luis_version":     internal.ResourceVersion(),
		},
	}
}

func configureFunc(d *schema.ResourceData) (interface{}, error) {
	options := clients.ClientOptions{
		Key:    d.Get("key").(string),
		Domain: d.Get("domain").(string),
	}

	return clients.NewClient(&options), nil
}

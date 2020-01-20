package luis

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

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
		DataSourcesMap: map[string]*schema.Resource{
			"luis_application": dataSourceApplication(),
			"luis_version":     dataSourceVersion(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"template_file": schema.DataSourceResourceShim(
				"template_file",
				dataSourceFile(),
			),
			"template_cloudinit_config": schema.DataSourceResourceShim(
				"template_cloudinit_config",
				dataSourceCloudinitConfig(),
			),
			"template_dir": resourceDir(),
		},
	}
}

func configureFunc(d *schema.ResourceData) (interface{}, error) {
	authoringKey := d.Get("authoring_key").(string)
	domain := d.Get("domain").(string)

	return client, nil
}

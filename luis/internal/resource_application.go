package internal

import (
	"fmt"

	luis "github.com/crazedpeanut/luis/client"
	"github.com/crazedpeanut/luis/client/operations"
	"github.com/crazedpeanut/luis/models"
	"github.com/hashicorp/terraform/helper/schema"
)

// ResourceApplication managed luis application resources
func ResourceApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Delete: resourceApplicationDelete,
		Update: resourceApplicationUpdate,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"culture": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"usage_scenario": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"initial_version_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.0.1",
			},
		},
	}
}

func resourceApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*luis.LuisAuthoring)

	id := d.Id()

	params := operations.GetApplicationParams{
		AppID: id,
	}
	resp, err := client.Operations.GetApplication(&params, nil)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(resp.Payload.ID)

	return nil
}

func resourceApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*luis.LuisAuthoring)

	application := models.ApplicationCreateObject{
		Name:             d.Get("name").(string),
		Description:      d.Get("desciption").(string),
		Culture:          d.Get("culture").(string),
		UsageScenario:    d.Get("usage_scenario").(string),
		Domain:           d.Get("domain").(string),
	}

	params := operations.CreateApplicationParams{
		ApplicationCreateObject: &application,
	}

	resp, err := client.Operations.CreateApplication(&params, nil)
	if err != nil {
		fmt.Errorf("Could not create application %s", err)
		return nil
	}

	d.SetId(resp.Payload)

	return nil
}

func resourceApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*luis.LuisAuthoring)

	application := models.ApplicationUpdateObject{
		Name:             d.Get("name").(string),
		Description:      d.Get("desciption").(string),
		UsageScenario:    d.Get("usage_scenario").(string),
		Domain:           d.Get("domain").(string),
	}

	params := operations.UpdateApplicationParams{
		ApplicationUpdateObject: &application,
		AppID: d.Id()
	}

	resp, err := client.Operations.UpdateApplication(&params, nil)
	if err != nil {
		fmt.Errorf("Could not update application %s", err)
		return nil
	}

	return nil
}

func resourceApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*luis.LuisAuthoring)

	id := d.Id()
	params := operations.DeleteApplicationParams{
		AppID: id,
	}

	d.SetId("")

	_, err := client.Operations.DeleteApplication(&params, nil)
	if err != nil {
		fmt.Errorf("Could not delete application %s", err)
	}

	return nil
}

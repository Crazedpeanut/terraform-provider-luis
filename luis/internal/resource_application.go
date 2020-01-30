package internal

import (
	"fmt"

	luis "github.com/crazedpeanut/go-luis-authoring-client/client"
	"github.com/crazedpeanut/go-luis-authoring-client/client/operations"
	"github.com/crazedpeanut/go-luis-authoring-client/models"
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
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
	client := meta.(*luis.Luis)

	id := d.Id()

	params := operations.NewGetApplicationParams()
	params.SetAppID(id)

	resp, err := client.Operations.GetApplication(params, nil)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error reading application %s %+v", id, err)
	}

	d.SetId(resp.Payload.ID)

	return nil
}

func resourceApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*luis.Luis)

	application := models.ApplicationCreateObject{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Culture:          d.Get("culture").(string),
		UsageScenario:    d.Get("usage_scenario").(string),
		Domain:           d.Get("domain").(string),
		InitialVersionID: d.Get("initial_version_id").(string),
	}

	params := operations.NewCreateApplicationParams()
	params.SetApplicationCreateObject(&application)

	resp, err := client.Operations.CreateApplication(params, nil)
	if err != nil {
		return fmt.Errorf("Error creating application %+v", err)
	}

	d.SetId(resp.Payload)

	return nil
}

func resourceApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*luis.Luis)
	id := d.Id()

	application := models.ApplicationUpdateObject{
		Name:        d.Get("name").(string),
		Description: d.Get("desciption").(string),
	}

	params := operations.NewUpdateApplicationParams()
	params.SetAppID(id)
	params.SetApplicationUpdateObject(&application)

	_, err := client.Operations.UpdateApplication(params, nil)
	if err != nil {
		return fmt.Errorf("Error updating application %s %+v", id, err)
	}

	return nil
}

func resourceApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*luis.Luis)

	id := d.Id()
	params := operations.NewDeleteApplicationParams()
	params.SetAppID(id)

	d.SetId("")

	_, err := client.Operations.DeleteApplication(params, nil)
	if err != nil {
		return fmt.Errorf("Error deleting application %s %+v", id, err)
	}

	return nil
}

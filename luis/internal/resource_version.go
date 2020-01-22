package internal

import (
	"encoding/json"
	"fmt"
	"time"

	luis "github.com/crazedpeanut/go-luis-authoring-client/client"
	"github.com/crazedpeanut/go-luis-authoring-client/client/operations"
	"github.com/crazedpeanut/go-luis-authoring-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

// ResourceVersion managed luis application resources
func ResourceVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceVersionCreate,
		Read:   resourceVersionRead,
		Delete: resourceVersionDelete,

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"trained": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"published": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*luis.LuisAuthoring)

	id := d.Id()
	appID := d.Get("app_id").(string)

	params := operations.NewGetApplicationVersionParams()
	params.SetAppID(appID)
	params.SetVersionID(id)

	_, err := client.Operations.GetApplicationVersion(params, nil)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error reading version %+v", err)
	}

	d.SetId(id)

	return nil
}

func resourceVersionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*luis.LuisAuthoring)

	appID := d.Get("app_id").(string)
	versionID := d.Get("version_id").(string)
	content := d.Get("content").(string)
	trained := d.Get("trained").(bool)
	published := d.Get("published").(bool)

	if isJSON(&content) {
		params := operations.NewImportVersionJSONParams()
		params.SetAppID(appID)
		params.SetVersionID(&versionID)
		params.SetContent(content)

		_, err := client.Operations.ImportVersionJSON(params, nil)
		if err != nil {
			return fmt.Errorf("Could not create version %+v", err)
		}
	} else {
		params := operations.NewImportVersionLuParams()
		params.SetAppID(appID)
		params.SetVersionID(&versionID)
		params.SetContent(content)

		_, err := client.Operations.ImportVersionLu(params, nil)
		if err != nil {
			return fmt.Errorf("Could not create version %+v", err)
		}
	}

	d.SetId(versionID)

	if trained {
		err := trainVersion(appID, versionID, client)
		if err != nil {
			return fmt.Errorf("Could not train version %+v", err)
		}
	}

	if published {
		err := publishVersion(appID, versionID, client)
		if err != nil {
			return fmt.Errorf("Could not publish version %+v", err)
		}
	}

	return nil
}

func resourceVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*luis.LuisAuthoring)

	appID := d.Get("app_id").(string)
	versionID := d.Get("version_id").(string)

	params := operations.NewDeleteVersionParams()
	params.SetAppID(appID)
	params.SetVersionID(versionID)

	d.SetId("")

	_, err := client.Operations.DeleteVersion(params, nil)
	if err != nil {
		return fmt.Errorf("Could not delete version %+v", err)
	}

	return nil
}

func isJSON(s *string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(*s), &js) == nil

}

func readJSONApp(raw *string) *models.JSONApp {
	var jsonApp models.JSONApp
	json.Unmarshal([]byte(*raw), &jsonApp)
	return &jsonApp
}

func publishVersion(appID string, versionID string, client *luis.LuisAuthoring) error {
	params := operations.NewPublishApplicationParams()
	params.SetAppID(appID)

	publishObject := models.ApplicationPublishObject{
		VersionID:            versionID,
		DirectVersionPublish: true,
	}
	params.SetApplicationPublishObject(&publishObject)

	_, _, err := client.Operations.PublishApplication(params, nil)

	return err
}

func trainVersion(appID string, versionID string, client *luis.LuisAuthoring) error {
	err := beginTrain(appID, versionID, client)
	if err != nil {
		return fmt.Errorf("Unable to start train %+v", err)
	}

	for {
		complete, err := isTrainComplete(appID, versionID, client)
		if err != nil {
			return fmt.Errorf("Unable to fetch train status %+v", err)
		}

		if complete {
			return nil
		}

		time.Sleep(10 * time.Second)
	}
}

func beginTrain(appID string, versionID string, client *luis.LuisAuthoring) error {
	params := operations.NewTrainVersionParams()
	params.SetAppID(appID)
	params.SetVersionID(versionID)

	time.Sleep(10 * time.Second)
	_, err := client.Operations.TrainVersion(params, nil)
	if err != nil {
		return err
	}
	return nil
}

func isTrainComplete(appID string, versionID string, client *luis.LuisAuthoring) (bool, error) {
	params := operations.NewGetApplicationVersionParams()
	params.SetAppID(appID)
	params.SetVersionID(versionID)

	resp, err := client.Operations.GetApplicationVersion(params, nil)
	if err != nil {
		return false, err
	}

	if resp.Payload.TrainingStatus == "Trained" {
		return true, nil
	}

	return false, nil
}

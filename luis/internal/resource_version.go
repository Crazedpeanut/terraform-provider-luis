package internal

import (
	"encoding/json"
	"fmt"
	"time"

	"log"

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
			"publish_version_direct": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"is_staging": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
		},
	}
}

func resourceVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*luis.Luis)

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
	client := meta.(*luis.Luis)

	appID := d.Get("app_id").(string)
	versionID := d.Get("version_id").(string)
	content := d.Get("content").(string)
	isStaging := d.Get("is_staging").(bool)
	trained := d.Get("trained").(bool)
	published := d.Get("published").(bool)
	directPublish := d.Get("publish_version_direct").(bool)

	if version, err := parseVersion(&content); err == nil {
		params := operations.NewImportVersionJSONParams()
		params.SetAppID(appID)
		params.SetVersionID(&versionID)
		params.SetJSONApp(&version)

		_, err := client.Operations.ImportVersionJSON(params, nil)
		if err != nil {
			if err, badRequest := err.(*operations.ImportVersionJSONBadRequest); badRequest {
				return fmt.Errorf("Could not create version bad request (JSON) %+v", err.GetPayload().Error)
			}
			return fmt.Errorf("Could not create version (JSON) %+v", err)
		}
	} else {
		log.Printf("[DEBUG] Not JSON %+v", err)
		params := operations.NewImportVersionLuParams()
		params.SetAppID(appID)
		params.SetVersionID(&versionID)
		params.SetContent(content)

		_, err := client.Operations.ImportVersionLu(params, nil)
		if err != nil {
			return fmt.Errorf("Could not create version (LuDown) %+v", err)
		}
	}

	d.SetId(versionID)

	if trained {
		err := trainVersion(appID, versionID, client)
		if err != nil {
			return fmt.Errorf("Could not train version %+v", err)
		}

		if published {
			err := publishVersion(appID, versionID, isStaging, directPublish, client)
			if err != nil {
				return fmt.Errorf("Could not publish version %+v", err)
			}
		}
	}

	return nil
}

func resourceVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*luis.Luis)

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

func parseVersion(s *string) (models.JSONApp, error) {
	var js models.JSONApp

	err := json.Unmarshal([]byte(*s), &js)

	return js, err
}

func readJSONApp(raw *string) *models.JSONApp {
	var jsonApp models.JSONApp
	json.Unmarshal([]byte(*raw), &jsonApp)
	return &jsonApp
}

func publishVersion(appID string, versionID string, isStaging bool, directPublish bool, client *luis.Luis) error {
	params := operations.NewPublishApplicationParams()
	params.SetAppID(appID)

	publishObject := models.ApplicationPublishObject{
		VersionID:            versionID,
		DirectVersionPublish: directPublish,
		IsStaging:            &isStaging,
	}
	params.SetApplicationPublishObject(&publishObject)

	_, _, err := client.Operations.PublishApplication(params, nil)

	return err
}

func trainVersion(appID string, versionID string, client *luis.Luis) error {
	log.Printf("[DEBUG] Begin training app: %s version: %s", appID, versionID)

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
			log.Printf("[DEBUG] training app: %s version: %s complete", appID, versionID)
			return nil
		}

		log.Printf("[DEBUG] training app: %s version: %s still going", appID, versionID)

		time.Sleep(10 * time.Second)
	}
}

func beginTrain(appID string, versionID string, client *luis.Luis) error {
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

func isTrainComplete(appID string, versionID string, client *luis.Luis) (bool, error) {
	params := operations.NewGetTrainingStatusParams()
	params.SetAppID(appID)
	params.SetVersionID(versionID)

	resp, err := client.Operations.GetTrainingStatus(params, nil)
	if err != nil {
		return false, err
	}

	log.Printf("[DEBUG] number of models training: %d", len(resp.Payload))

	for _, model := range resp.Payload {
		log.Printf("[DEBUG] Model status: %s", model.Details.Status)

		if model.Details.Status == "Fail" {
			return false, fmt.Errorf("Error training version %s", model.Details.FailureReason)
		}

		if model.Details.Status == "InProgress" || model.Details.Status == "Queued" {
			return false, nil
		}
	}

	return true, nil
}

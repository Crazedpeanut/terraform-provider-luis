package luis

import (
	"io/ioutil"
	"net/http"
)

type Application struct {
	name string
}

type CreateApplicationModel struct {
	name string
}

type UpdateApplicationModel struct {
	name string
}

type GetApplications func() []Application
type GetApplication func(string) Application
type DeleteApplication func(string)
type CreateApplication func(CreateApplicationModel) string
type UpdateApplication func(string, UpdateApplicationModel)

type Client struct {
	getApps   GetApplications
	getApp    GetApplication
	deleteApp DeleteApplication
	createApp CreateApplication
	updateApp UpdateApplication
}

type ClientOptions struct {
	authoringKey string
	domain       string
}

func createApplication(o *ClientOptions) string {
	return func(app *CreateApplication) {

	}
}

func createApplication(o *ClientOptions) string {
	return func(id string, app *UpdateApplication) {

	}
}

func getApplications(o *ClientOptions) string {
	return func() []Application {

	}
}

func getApplication(o *ClientOptions) string {
	return func(id string) []Application {

	}
}

func getApplication(o *ClientOptions) string {
	return func(id string) {

	}
}

func NewClient(o *ClientOptions) *Client {

	return &Client{
		createApp: createApplication(o),
		getApps: getApplications(o),
		getApp: getApplication(o)
		updateApp: updateApplication(o),
		deleteApp: deleteApplication(o)
	}
}

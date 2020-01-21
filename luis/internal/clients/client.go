package clients

import (
	"github.com/go-openapi/strfmt"

	luis "github.com/crazedpeanut/luis/client"
	"github.com/go-openapi/runtime/client"
)

// ClientOptions is used to configure luis clients
type ClientOptions struct {
	AuthoringKey string
	Domain       string
}

// NewClient produces a client to call LUIS Authoring API
func NewClient(o *ClientOptions) *luis.LuisAuthoring {

	transport := client.New(o.Domain, "", nil)
	transport.DefaultAuthentication = client.APIKeyAuth("Ocp-Apim-Subscription-Key", "header", o.AuthoringKey)

	return luis.New(transport, strfmt.Default)
}

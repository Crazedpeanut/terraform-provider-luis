package clients

import (
	"github.com/go-openapi/strfmt"

	luis "github.com/crazedpeanut/go-luis-authoring-client/client"
	"github.com/go-openapi/runtime/client"
)

// ClientOptions is used to configure luis clients
type ClientOptions struct {
	Key    string
	Domain string
}

// NewClient produces a client to call LUIS Authoring API
func NewClient(o *ClientOptions) *luis.Luis {

	transport := client.New(o.Domain, "/luis/api/v2.0", nil)
	transport.DefaultAuthentication = client.APIKeyAuth("Ocp-Apim-Subscription-Key", "header", o.Key)

	return luis.New(transport, strfmt.Default)
}

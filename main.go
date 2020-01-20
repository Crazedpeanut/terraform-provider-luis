package main

import (
	"github.com/crazedpeanut/terraform-provider-luis/template"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: template.Provider})
}

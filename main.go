package main

import (
	"github.com/crazedpeanut/terraform-provider-luis/luis"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: luis.Provider})
}

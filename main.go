package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/wleese/terraform-provider-calico/calico"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: calico.Provider,
	})
}

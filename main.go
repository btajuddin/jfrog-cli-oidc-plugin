package main

import (
	"github.com/btajuddin/jfrog-cli-oidc-plugin/commands"
	"github.com/jfrog/jfrog-cli-core/v2/plugins"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
)

func main() {
	plugins.PluginMain(getApp())
}

func getApp() components.App {
	app := components.App{}
	app.Name = "oidc-exchange"
	app.Description = "Easily retrieve an access token from and OIDC token."
	app.Version = "v1.0.2"
	app.Commands = getCommands()
	return app
}

func getCommands() []components.Command {
	return []components.Command{
		commands.GetExchangeCommand()}
}

package cmd

import (
	"log"

	"github.com/antham/wo/workspace"
	"github.com/spf13/viper"
)

var wksManager workspaceManager

func newWorkspaceManager() workspaceManager {
	if wksManager != nil {
		return wksManager
	}
	viper.AutomaticEnv()
	options := []func(*workspace.WorkspaceManager){
		workspace.WithEditor(viper.GetString("EDITOR"), viper.GetString("VISUAL")),
		workspace.WithShellPath(viper.GetString("SHELL")),
	}
	viper.SetEnvPrefix("WO")
	if viper.IsSet("CONFIG_PATH") {
		options = append(options, workspace.WithConfigPath(viper.GetString("CONFIG_PATH")))
	}
	var err error
	wksManager, err = workspace.NewWorkspaceManager(options...)
	if err != nil {
		log.Fatal(err)
	}
	return wksManager
}

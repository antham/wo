package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newGlobalGetCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	return &cobra.Command{
		Use:               "get key",
		Short:             "Get a configuration",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completionManager.Process,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "config-dir":
				cmd.Printf("%s", workspaceManager.GetConfigDir())
			default:
				return fmt.Errorf("Key '%s' does not exist", args[0])
			}
			return nil
		},
	}
}

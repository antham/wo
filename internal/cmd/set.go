package cmd

import (
	"github.com/spf13/cobra"
)

func newSetCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	return &cobra.Command{
		Use:               "set workspace key value",
		Short:             "Set a configuration",
		Args:              cobra.ExactArgs(3),
		ValidArgsFunction: completionManager.Process,
		RunE: func(cmd *cobra.Command, args []string) error {
			return workspaceManager.SetConfig(args[0], map[string]string{args[1]: args[2]})
		},
	}
}

package cmd

import (
	"github.com/spf13/cobra"
)

func newEditEnvCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	return &cobra.Command{
		Use:               "edit workspace environment",
		Short:             "Edit a workspace environment",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: completionManager.Process,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return workspaceManager.EditEnv(args[0], args[1])
		},
	}
}

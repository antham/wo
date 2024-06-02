package cmd

import (
	"github.com/spf13/cobra"
)

func newCreateCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	return &cobra.Command{
		Use:               "create workspace project-path",
		Short:             "Create a workspace",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: completionManager.Process,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return workspaceManager.Create(args[0], args[1])
		},
	}
}

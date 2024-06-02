package cmd

import (
	"github.com/spf13/cobra"
)

func newRemoveCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	return &cobra.Command{
		Use:               "remove workspace",
		Short:             "Remove a workspace",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completionManager.Process,
		RunE: func(cmd *cobra.Command, args []string) error {
			return workspaceManager.Remove(args[0])
		},
	}
}

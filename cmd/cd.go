package cmd

import (
	"github.com/spf13/cobra"
)

func newCdCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	return &cobra.Command{
		Use:               "cd",
		Short:             "Jump to the workspace directory",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completionManager.Process,
		RunE: func(cmd *cobra.Command, args []string) error {
			return workspaceManager.Cd(args[0])
		},
	}
}

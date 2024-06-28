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
			err := workspaceManager.Remove(args[0])
			if err != nil {
				return err
			}
			cmd.Printf("Workspace '%s' deleted\n", args[0])
			return nil
		},
	}
}

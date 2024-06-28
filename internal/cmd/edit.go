package cmd

import (
	"github.com/spf13/cobra"
)

func newEditCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	return &cobra.Command{
		Use:               "edit workspace",
		Short:             "Edit a workspace",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completionManager.Process,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := workspaceManager.Edit(args[0])
			if err != nil {
				return err
			}
			cmd.Printf("Workspace '%s' edited\n", args[0])
			return nil
		},
	}
}

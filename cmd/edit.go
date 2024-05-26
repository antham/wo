package cmd

import (
	"github.com/spf13/cobra"
)

func newEditCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	return &cobra.Command{
		Use:               "edit workspace [environment]",
		Short:             "Edit a workspace",
		Args:              cobra.RangeArgs(1, 2),
		ValidArgsFunction: completionManager.Process,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			switch len(args) {
			case 1:
				err = workspaceManager.Edit(args[0])
			case 2:
				err = workspaceManager.EditEnv(args[0], args[1])
			}
			return err
		},
	}
}

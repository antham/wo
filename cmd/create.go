package cmd

import (
	"github.com/spf13/cobra"
)

func newCreateCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	return &cobra.Command{
		Use:               "create workspace [environment]",
		Short:             "Create a workspace",
		Args:              cobra.RangeArgs(1, 2),
		ValidArgsFunction: completionManager.Process,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			switch len(args) {
			case 1:
				err = workspaceManager.Create(args[0])
			case 2:
				err = workspaceManager.CreateEnv(args[0], args[1])
			}
			return err
		},
	}
}

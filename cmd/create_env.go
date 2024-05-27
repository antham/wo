package cmd

import (
	"github.com/spf13/cobra"
)

func newCreateEnvCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	return &cobra.Command{
		Use:               "create workspace environment",
		Short:             "Create a workspace environment",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: completionManager.Process,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return workspaceManager.CreateEnv(args[0], args[1])
		},
	}
}

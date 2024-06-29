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
		RunE: func(cmd *cobra.Command, args []string) error {
			err := workspaceManager.EditEnv(args[0], args[1])
			if err != nil {
				return err
			}
			cmd.Printf(
				regularStyle.Render("Environment '")+highlightedStyle.Render("%s")+regularStyle.Render("' edited on workspace '")+highlightedStyle.Render("%s")+regularStyle.Render("'")+"\n",
				args[1], args[0],
			)
			return nil
		},
	}
}

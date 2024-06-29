package cmd

import (
	"errors"

	"github.com/antham/wo/internal/cmd/internal/validator"
	"github.com/spf13/cobra"
)

func newCreateEnvCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	return &cobra.Command{
		Use:               "create workspace environment",
		Short:             "Create a workspace environment",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: completionManager.Process,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return errors.Join(validator.ValidateName(args[1]))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := workspaceManager.CreateEnv(args[0], args[1])
			if err != nil {
				return err
			}
			cmd.Printf(
				regularStyle.Render("Environment '")+highlightedStyle.Render("%s")+regularStyle.Render("' added on workspace '")+highlightedStyle.Render("%s")+regularStyle.Render("'")+"\n",
				args[1], args[0],
			)
			return nil
		},
	}
}

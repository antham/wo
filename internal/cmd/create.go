package cmd

import (
	"errors"

	"github.com/antham/wo/internal/cmd/internal/validator"
	"github.com/spf13/cobra"
)

func newCreateCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	return &cobra.Command{
		Use:               "create workspace project-path",
		Short:             "Create a workspace",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: completionManager.Process,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return errors.Join(validator.ValidateName(args[0]))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := workspaceManager.Create(args[0], args[1])
			if err != nil {
				return err
			}
			cmd.Printf(regularStyle.Render("Workspace '")+highlightedStyle.Render("%s")+regularStyle.Render("' created on path '")+highlightedStyle.Render("%s")+regularStyle.Render("', reload your shell to setup the project aliases")+"\n", args[0], args[1])
			return nil
		},
	}
}

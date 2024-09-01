package cmd

import (
	"github.com/spf13/cobra"
)

func newFixCmd(workspaceManager workspaceManager) *cobra.Command {
	return &cobra.Command{
		Use:   "fix",
		Short: "Fix the possible failures in the config folder",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := workspaceManager.Fix()
			if err != nil {
				return err
			}
			cmd.Print(regularStyle.Render("Config folder fixed"))
			return nil
		},
	}
}

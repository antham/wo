package cmd

import (
	"github.com/spf13/cobra"
)

func newMigrateCmd(workspaceManager workspaceManager) *cobra.Command {
	return &cobra.Command{
		Use:    "migrate",
		Short:  "Migrate the existing config",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := workspaceManager.Migrate()
			if err != nil {
				return err
			}
			cmd.Print(regularStyle.Render("Config migrated"))
			return nil
		},
	}
}

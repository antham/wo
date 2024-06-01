package cmd

import (
	"github.com/spf13/cobra"
)

func newAliasCmd(workspaceManager workspaceManager) *cobra.Command {
	var prefix string
	cmd := &cobra.Command{
		Use:   "alias",
		Short: "List all shell aliases",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			aliases, err := workspaceManager.BuildAliases(prefix)
			if err != nil {
				return err
			}
			for _, alias := range aliases {
				cmd.Println(alias)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&prefix, "prefix", "p", "", "Prefix name to use for the aliases")
	return cmd
}

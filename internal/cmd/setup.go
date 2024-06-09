package cmd

import (
	"github.com/spf13/cobra"
)

func newSetupCmd(workspaceManager workspaceManager) *cobra.Command {
	var prefix string
	cmd := &cobra.Command{
		Use:   "setup shell",
		Short: "Setup wo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := &cobra.Command{
				Use: "wo",
			}
			var err error
			switch args[0] {
			case "bash":
				err = c.GenBashCompletionV2(cmd.OutOrStdout(), true)
			case "fish":
				err = c.GenFishCompletion(cmd.OutOrStdout(), true)
			case "zsh":
				err = c.GenZshCompletion(cmd.OutOrStdout())
			}
			if err != nil {
				return err
			}

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

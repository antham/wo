package cmd

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
)

func newSetupCmd(workspaceManager workspaceManager) *cobra.Command {
	var prefix string
	cmd := &cobra.Command{
		Use:       "setup shell",
		Short:     "Command to setup wo in the shell",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "fish", "zsh", "sh"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !slices.Contains(cmd.ValidArgs, args[0]) {
				return fmt.Errorf("the first argument must one of among: %v", cmd.ValidArgs)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// We need this to be able to have the completion working
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
	cmd.Flags().StringVarP(&prefix, "prefix", "p", "c_", "Prefix name to use for the aliases")
	return cmd
}

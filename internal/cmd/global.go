package cmd

import "github.com/spf13/cobra"

func newGlobalCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "global",
		Short: "Manage global features",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}

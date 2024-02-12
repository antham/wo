package cmd

import (
	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove workspace",
	Short: "Remove a workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := workspace.NewWorkspaceManager()
		if err != nil {
			return err
		}
		return w.Remove(args[0])
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

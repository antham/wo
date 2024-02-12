package cmd

import (
	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit workspace",
	Short: "Edit a workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := workspace.NewWorkspaceManager()
		if err != nil {
			return err
		}
		return w.Edit(args[0])
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}

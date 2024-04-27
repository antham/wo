package cmd

import (
	"github.com/antham/wo/cmd/internal/completion"
	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove workspace",
	Short: "Remove a workspace",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		w, err := workspace.NewWorkspaceManager()
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}
		c := completion.New(w)
		switch len(args) {
		case 0:
			workspaces, err := c.FindWorkspaces(toComplete)
			if err != nil {
				return []string{}, cobra.ShellCompDirectiveNoFileComp
			}
			return workspaces, cobra.ShellCompDirectiveNoFileComp
		}
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	},
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

package cmd

import (
	"github.com/antham/wo/cmd/internal/completion"
	"github.com/spf13/cobra"
)

func newRemoveCmd(workspaceManager workspaceManager) *cobra.Command {
	return &cobra.Command{
		Use:   "remove workspace",
		Short: "Remove a workspace",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			c := completion.New(workspaceManager)
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
			return workspaceManager.Remove(args[0])
		},
	}
}

func init() {
	rootCmd.AddCommand(newRemoveCmd(newWorkspaceManager()))
}

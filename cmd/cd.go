package cmd

import (
	"github.com/antham/wo/cmd/internal/completion"
	"github.com/spf13/cobra"
)

func newCdCmd(workspaceManager workspaceManager) *cobra.Command {
	return &cobra.Command{
		Use:   "cd",
		Short: "Jump to the workspace directory",
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
			return workspaceManager.Cd(args[0])
		},
	}
}

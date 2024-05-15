package cmd

import (
	"github.com/antham/wo/cmd/internal/completion"
	"github.com/spf13/cobra"
)

func newSetCmd(workspaceManager workspaceManager) *cobra.Command {
	return &cobra.Command{
		Use:   "set workspace key value",
		Short: "Set a configuration",
		Args:  cobra.ExactArgs(3),
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
			return workspaceManager.SetConfig(args[0], args[1], args[2])
		},
	}
}

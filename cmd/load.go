package cmd

import (
	"github.com/antham/wo/cmd/internal/completion"
	"github.com/spf13/cobra"
)

var env string

func newLoadCmd(workspaceManager workspaceManager) *cobra.Command {
	loadCmd := &cobra.Command{
		Use:     "load workspace [environment]",
		Aliases: []string{"l"},
		Short:   "Load a workspace",
		Args:    cobra.RangeArgs(1, 2),
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
			switch len(args) {
			case 1:
				return workspaceManager.Load(args[0], "")
			case 2:
				return workspaceManager.Load(args[0], args[1])
			}
			return nil
		},
	}
	return loadCmd
}

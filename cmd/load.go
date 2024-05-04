package cmd

import (
	"github.com/antham/wo/cmd/internal/completion"
	"github.com/spf13/cobra"
)

var env string

func newLoadCmd(workspaceManager workspaceManager) *cobra.Command {
	loadCmd := &cobra.Command{
		Use:     "load workspace",
		Aliases: []string{"l"},
		Short:   "Load a workspace",
		Args:    cobra.ExactArgs(1),
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
			return workspaceManager.Load(args[0], env)
		},
	}
	loadCmd.Flags().StringVarP(&env, "env", "e", "", "Environment to use (e.g. prod)")
	return loadCmd
}

package cmd

import (
	"github.com/antham/wo/cmd/internal/completion"
	"github.com/spf13/cobra"
)

func newEditCmd(workspaceManager workspaceManager) *cobra.Command {
	return &cobra.Command{
		Use:   "edit workspace [environment]",
		Short: "Edit a workspace",
		Args:  cobra.RangeArgs(1, 2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			c := completion.New(workspaceManager)
			switch len(args) {
			case 0:
				workspaces, err := c.FindWorkspaces(toComplete)
				if err != nil {
					return []string{}, cobra.ShellCompDirectiveNoFileComp
				}
				return workspaces, cobra.ShellCompDirectiveNoFileComp
			case 1:
				envs, err := c.FindEnvs(args[0], toComplete)
				if err != nil {
					return []string{}, cobra.ShellCompDirectiveNoFileComp
				}
				return envs, cobra.ShellCompDirectiveNoFileComp
			}
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			switch len(args) {
			case 1:
				err = workspaceManager.Edit(args[0])
			case 2:
				err = workspaceManager.EditEnv(args[0], args[1])
			}
			return err
		},
	}
}

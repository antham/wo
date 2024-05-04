/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/antham/wo/cmd/internal/completion"
	"github.com/spf13/cobra"
)

func newCreateCmd(workspaceManager workspaceManager) *cobra.Command {
	return &cobra.Command{
		Use:   "create workspace [environment]",
		Short: "Create a workspace",
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
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			switch len(args) {
			case 1:
				err = workspaceManager.Create(args[0])
			case 2:
				err = workspaceManager.CreateEnv(args[0], args[1])
			}
			return err
		},
	}
}

func init() {
	rootCmd.AddCommand(newCreateCmd(newWorkspaceManager()))
}

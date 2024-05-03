package cmd

import (
	"os"
	"os/exec"

	"github.com/antham/wo/cmd/internal/completion"
	"github.com/spf13/cobra"
)

func newRunCmd(workspaceManager workspaceManager) *cobra.Command {
	return &cobra.Command{
		Use:     "run workspace function [function-args]...",
		Aliases: []string{"r"},
		Short:   "Run a function in a given workspace",
		Args:    cobra.MinimumNArgs(2),
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
				commands, err := c.FindFunctions(args[0], toComplete)
				if err != nil {
					return []string{}, cobra.ShellCompDirectiveNoFileComp
				}
				return commands, cobra.ShellCompDirectiveNoFileComp
			}
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := workspaceManager.RunFunction(args[0], env, args[1:])
			if exitError, ok := err.(*exec.ExitError); ok {
				os.Exit(exitError.ExitCode())
			}
			return nil
		},
	}
}

func init() {
	runCmd := newRunCmd(newWorkspaceManager())
	runCmd.Flags().StringVarP(&env, "env", "e", "", "Environment to use (e.g. prod)")
	rootCmd.AddCommand(runCmd)
}

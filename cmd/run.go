package cmd

import (
	"github.com/antham/wo/cmd/internal/completion"
	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:     "run workspace function [function-args]...",
	Aliases: []string{"r"},
	Short:   "Run a function in a given workspace",
	Args:    cobra.MinimumNArgs(2),
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
		w, err := workspace.NewWorkspaceManager()
		if err != nil {
			return err
		}
		return w.RunFunction(args[0], env, args[1:])
	},
}

func init() {
	runCmd.Flags().StringVarP(&env, "env", "e", "", "Environment to use (e.g. prod)")
	rootCmd.AddCommand(runCmd)
}

package cmd

import (
	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:     "run workspace function [function-args]...",
	Aliases: []string{"r"},
	Short:   "Run a command in a given workspace",
	Args:    cobra.MinimumNArgs(2),
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

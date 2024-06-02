package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newRunCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	runCmd := &cobra.Command{
		Use:               "run workspace function [function-args]...",
		Aliases:           []string{"r"},
		Short:             "Run a function in a given workspace",
		Args:              cobra.MinimumNArgs(2),
		ValidArgsFunction: completionManager.Process,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := workspaceManager.RunFunction(args[0], env, args[1:])
			if exitError, ok := err.(*exec.ExitError); ok {
				os.Exit(exitError.ExitCode())
			}
			return err
		},
	}
	runCmd.Flags().StringVarP(&env, "env", "e", "", "Environment to use (e.g. prod)")
	return runCmd
}

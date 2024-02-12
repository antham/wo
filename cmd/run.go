package cmd

import (
	"log"

	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run workspace function [function-args]...",
	Short: "Run a command in a given workspace",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		w, err := workspace.NewWorkspaceManager()
		if err != nil {
			log.Fatal(err)
		}
		err = w.RunFunction(args[0], env, args[1:])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	runCmd.Flags().StringVarP(&env, "env", "e", "", "Environment to use (e.g. prod)")
	rootCmd.AddCommand(runCmd)
}

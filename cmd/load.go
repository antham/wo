package cmd

import (
	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

var env string

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:     "load workspace",
	Aliases: []string{"l"},
	Short:   "Load a workspace",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := workspace.NewWorkspaceManager()
		if err != nil {
			return err
		}
		return w.Load(args[0], env)
	},
}

func init() {
	loadCmd.Flags().StringVarP(&env, "env", "e", "", "Environment to use (e.g. prod)")
	rootCmd.AddCommand(loadCmd)
}

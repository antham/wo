package cmd

import (
	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

// editEnvCmd represents the editEnv command
var editEnvCmd = &cobra.Command{
	Use:   "edit-env workspace [environment]",
	Short: "Edit an environment for a given workspace",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := workspace.NewWorkspaceManager()
		if err != nil {
			return err
		}
		env := ""
		if len(args) > 1 {
			env = args[1]
		}
		return w.EditEnv(args[0], env)
	},
}

func init() {
	rootCmd.AddCommand(editEnvCmd)
}

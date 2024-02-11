/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

// editEnvCmd represents the editEnv command
var editEnvCmd = &cobra.Command{
	Use:   "edit-env workspace [environment]",
	Short: "Edit an environment for a given workspace",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		w, err := workspace.NewWorkspaceManager()
		if err != nil {
			log.Fatal(err)
		}
		env := ""
		if len(args) > 1 {
			env = args[1]
		}
		err = w.EditEnv(args[0], env)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(editEnvCmd)
}

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

var env string

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load workspace",
	Short: "Load a workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		w, err := workspace.NewWorkspaceManager()
		if err != nil {
			log.Fatal(err)
		}
		err = w.Load(args[0], env)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	loadCmd.Flags().StringVarP(&env, "env", "e", "", "Environment to use (e.g. prod)")
	rootCmd.AddCommand(loadCmd)
}

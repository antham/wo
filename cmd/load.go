package cmd

import (
	"log"

	"github.com/antham/wo/cmd/internal/completion"
	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

var env string

func newloadCmd(workspaceManager workspaceManager) *cobra.Command {
	return &cobra.Command{
		Use:     "load workspace",
		Aliases: []string{"l"},
		Short:   "Load a workspace",
		Args:    cobra.ExactArgs(1),
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return workspaceManager.Load(args[0], env)
		},
	}
}

func init() {
	w, err := workspace.NewWorkspaceManager()
	if err != nil {
		log.Fatal(err)
	}
	loadCmd := newloadCmd(w)
	loadCmd.Flags().StringVarP(&env, "env", "e", "", "Environment to use (e.g. prod)")
	rootCmd.AddCommand(loadCmd)
}

package cmd

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "list",
	Short: "List workspace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := workspace.NewWorkspaceManager()
		if err != nil {
			return err
		}
		workspaces, err := s.List()
		if err != nil {
			return err
		}
		slices.SortFunc(workspaces, func(a, b workspace.Workspace) int {
			return cmp.Compare(a.Name, b.Name)
		})
		for _, w := range workspaces {
			fmt.Println(w.Name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

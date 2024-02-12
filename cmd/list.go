package cmd

import (
	"fmt"
	"sort"

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
		sort.Sort(workspace.ByName(workspaces))
		for _, w := range workspaces {
			fmt.Println(w.Name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

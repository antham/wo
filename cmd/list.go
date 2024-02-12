package cmd

import (
	"fmt"
	"log"
	"sort"

	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "list",
	Short: "List workspace",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := workspace.NewWorkspaceManager()
		if err != nil {
			log.Fatal(err)
		}
		workspaces := s.List()
		if err != nil {
			log.Fatal(err)
		}
		sort.Sort(workspace.ByName(workspaces))
		for _, w := range workspaces {
			fmt.Println(w.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

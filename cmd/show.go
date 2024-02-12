package cmd

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show workspace",
	Short: "Show all functions available in a workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		w, err := workspace.NewWorkspaceManager()
		if err != nil {
			log.Fatal(err)
		}
		wo, err := w.Get(args[0])
		if err != nil {
			log.Fatal(err)
		}
		l := 0
		for _, c := range wo.Commands {
			if len(c.Command) > l {
				l = len(c.Command)
			}
		}
		sort.Sort(workspace.ByCommand(wo.Commands))
		fmt.Println(wo.Name)
		if len(wo.Commands) == 0 {
			fmt.Println("   no functions defined")
		}
		for _, c := range wo.Commands {
			fmt.Printf("   %-"+strconv.Itoa(l)+"s", c.Command)
			if c.Description != "" {
				fmt.Printf(" - %s\n", c.Description)
			} else {
				fmt.Println()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}

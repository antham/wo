package cmd

import (
	"fmt"

	"github.com/antham/wo/workspace"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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
		wss := []string{}
		workspaceRowTableSize := 11
		for _, w := range workspaces {
			if len(w.Name)+1 > workspaceRowTableSize {
				workspaceRowTableSize = len(w.Name) + 1
			}
			wss = append(wss, w.Name)
		}
		ws := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#C683D7"))).
			Headers("Workspaces").
			StyleFunc(func(row, col int) lipgloss.Style {
				var style lipgloss.Style
				switch {
				case row == 0:
					style = style.Bold(true).Foreground(lipgloss.Color("#C683D7"))
				}
				style = style.Copy().Width(workspaceRowTableSize)
				return style
			}).
			Rows([][]string{wss}...)
		fmt.Println(ws)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

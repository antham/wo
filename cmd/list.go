package cmd

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newListCmd(workspaceManager workspaceManager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List workspace",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaces, err := workspaceManager.List()
			if err != nil {
				return err
			}
			if len(workspaces) == 0 {
				return errors.New("no workspaces defined")
			}
			separator := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#22668D")).
				Render("---")
			title := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFCC70")).
				Render("Workspaces")
			var list string
			for _, w := range workspaces {
				list += lipgloss.NewStyle().
					Foreground(lipgloss.Color("#8ECDDD")).
					Render("• " + w.Name)
			}
			fmt.Println(title)
			fmt.Println(separator)
			fmt.Println(list)
			return nil
		},
	}
}

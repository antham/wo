package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newListCmd(workspaceManager workspaceManager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List workspaces",
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
			var list []string
			for _, w := range workspaces {
				list = append(list, lipgloss.NewStyle().
					Foreground(lipgloss.Color("#8ECDDD")).
					Render(fmt.Sprintf("* %s", w.Name)))
			}
			cmd.Println(title)
			cmd.Println(separator)
			cmd.Println(strings.Join(list, "\n"))
			return nil
		},
	}
}

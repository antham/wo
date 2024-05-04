package cmd

import (
	"fmt"
	"strings"

	"github.com/antham/wo/cmd/internal/completion"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newShowCmd(workspaceManager workspaceManager) *cobra.Command {
	return &cobra.Command{
		Use:   "show workspace",
		Short: "Show all functions available in a workspace",
		Args:  cobra.ExactArgs(1),
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
			wo, err := workspaceManager.Get(args[0])
			if err != nil {
				return err
			}
			separator := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#22668D")).
				Render("---")
			title := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFCC70")).
				Render(fmt.Sprintf("Workspace %s", wo.Name))
			functionTitle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFCC70")).
				Height(2).
				Render("Functions")
			var functions []string
			for _, f := range wo.Functions {
				functions = append(
					functions,
					fmt.Sprintf(
						"%s %s %s",
						lipgloss.NewStyle().
							Foreground(lipgloss.Color("#8ECDDD")).
							Render("•"),
						lipgloss.NewStyle().
							Foreground(lipgloss.Color("#FFFADD")).
							Render(f.Function),
						lipgloss.NewStyle().
							Foreground(lipgloss.Color("#8ECDDD")).
							Render(fmt.Sprintf(": %s", f.Description)),
					),
				)
			}
			if len(wo.Functions) == 0 {
				functions = append(functions, lipgloss.NewStyle().
					Foreground(lipgloss.Color("#8ECDDD")).
					Render("No functions defined"))
			}
			envTitle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFCC70")).
				Height(2).
				Render("Envs")
			var envs []string
			for _, e := range wo.Envs {
				envs = append(envs, lipgloss.NewStyle().
					Foreground(lipgloss.Color("#8ECDDD")).
					Render(fmt.Sprintf("• %s", e)))
			}
			if len(wo.Envs) == 0 {
				envs = append(envs, lipgloss.NewStyle().
					Foreground(lipgloss.Color("#8ECDDD")).
					Render("No envs defined"))
			}
			fmt.Println(title)
			fmt.Println(separator)
			fmt.Println(functionTitle)
			fmt.Println(strings.Join(functions, "\n"))
			fmt.Println(separator)
			fmt.Println(envTitle)
			fmt.Println(strings.Join(envs, "\n"))
			return nil
		},
	}
}

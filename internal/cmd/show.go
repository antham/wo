package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newShowCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	return &cobra.Command{
		Use:               "show workspace",
		Short:             "Show functions and envs available in a workspace",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completionManager.Process,
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
				Render("Functions")
			var functions []string
			for _, f := range wo.Functions.Functions {
				description := ""
				if f.Description != "" {
					description = lipgloss.NewStyle().
						Foreground(lipgloss.Color("#8ECDDD")).
						Render(fmt.Sprintf(" : %s", f.Description))
				}
				functions = append(
					functions,
					fmt.Sprintf(
						"%s %s%s",
						lipgloss.NewStyle().
							Foreground(lipgloss.Color("#8ECDDD")).
							Render("*"),
						lipgloss.NewStyle().
							Foreground(lipgloss.Color("#FFFADD")).
							Render(f.Name),
						description,
					),
				)
			}
			if len(wo.Functions.Functions) == 0 {
				functions = append(functions, lipgloss.NewStyle().
					Foreground(lipgloss.Color("#8ECDDD")).
					Render("No functions defined"))
			}
			envTitle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFCC70")).
				Render("Envs")
			var envs []string
			for _, e := range wo.Envs {
				envs = append(envs, lipgloss.NewStyle().
					Foreground(lipgloss.Color("#8ECDDD")).
					Render(fmt.Sprintf("* %s", e.Name)))
			}
			if len(wo.Envs) == 0 {
				envs = append(envs, lipgloss.NewStyle().
					Foreground(lipgloss.Color("#8ECDDD")).
					Render("No envs defined"))
			}
			cmd.Println(title)
			cmd.Println(separator)
			cmd.Println(functionTitle)
			cmd.Println()
			cmd.Println(strings.Join(functions, "\n"))
			cmd.Println(separator)
			cmd.Println(envTitle)
			cmd.Println()
			cmd.Println(strings.Join(envs, "\n"))
			return nil
		},
	}
}

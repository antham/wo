package cmd

import (
	"fmt"

	"github.com/antham/wo/cmd/internal/completion"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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
			functionRowTableSize := []int{10, 12}
			fs := [][]string{}
			for _, c := range wo.Functions {
				if c.Description == "" {
					c.Description = "-"
				}
				fs = append(fs, []string{c.Function, c.Description})
				if len(c.Function)+1 > functionRowTableSize[0] {
					functionRowTableSize[0] = len(c.Function) + 1
				}
				if len(c.Description)+1 > functionRowTableSize[1] {
					functionRowTableSize[1] = len(c.Description) + 1
				}
			}
			envRowTableSize := 5
			for _, e := range wo.Envs {
				if len(e)+1 > envRowTableSize {
					envRowTableSize = len(e) + 1
				}
			}
			title := lipgloss.NewStyle().
				MarginBottom(1).
				Bold(true).
				Foreground(lipgloss.Color("#7071E8")).
				Render(fmt.Sprintf("Workspace %s", wo.Name))
			functions := table.New().
				Border(lipgloss.NormalBorder()).
				BorderRow(true).
				BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#C683D7"))).
				Headers("Functions", "Description").
				StyleFunc(func(row, col int) lipgloss.Style {
					var style lipgloss.Style
					switch {
					case row == 0:
						style = style.Bold(true).Foreground(lipgloss.Color("#C683D7"))
					}
					style = style.Copy().Width(functionRowTableSize[col])
					return style
				}).
				Rows(fs...)
			envs := table.New().
				Border(lipgloss.NormalBorder()).
				BorderRow(true).
				BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#C683D7"))).
				Headers("Envs").
				StyleFunc(func(row, col int) lipgloss.Style {
					var style lipgloss.Style
					switch {
					case row == 0:
						style = style.Bold(true).Foreground(lipgloss.Color("#C683D7"))
					}
					style = style.Copy().Width(functionRowTableSize[col])
					return style
				}).
				Rows([][]string{wo.Envs}...)
			fmt.Println(title)
			fmt.Println(functions)
			fmt.Println(envs)
			return nil
		},
	}
}

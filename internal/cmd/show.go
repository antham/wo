package cmd

import (
	"fmt"
	"sort"
	"strings"

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
			title := titleStyle.
				Render(fmt.Sprintf("Workspace %s", wo.Name))
			configTitle := titleStyle.
				Render("Configuration")
			var configs []string
			for key, value := range wo.Config {
				configs = append(
					configs,
					fmt.Sprintf(
						"%s %s%s",
						regularStyle.
							Render("*"),
						highlightedStyle.
							Render(key),
						regularStyle.
							Render(fmt.Sprintf(" : %s", value)),
					),
				)
			}
			sort.Strings(configs)
			functionTitle := titleStyle.
				Render("Functions")
			var functions []string
			for _, f := range wo.Functions.Functions {
				description := ""
				if f.Description != "" {
					description = regularStyle.
						Render(fmt.Sprintf(" : %s", f.Description))
				}
				functions = append(
					functions,
					fmt.Sprintf(
						"%s %s%s",
						regularStyle.
							Render("*"),
						highlightedStyle.
							Render(f.Name),
						description,
					),
				)
			}
			if len(wo.Functions.Functions) == 0 {
				functions = append(functions, regularStyle.
					Render("No functions defined"))
			}
			envTitle := titleStyle.
				Render("Envs")
			var envs []string
			for _, e := range wo.Envs {
				envs = append(envs, regularStyle.
					Render(fmt.Sprintf("* %s", e.Name)))
			}
			if len(wo.Envs) == 0 {
				envs = append(envs, regularStyle.
					Render("No envs defined"))
			}
			cmd.Println(title)
			cmd.Println()
			cmd.Println(separator)
			cmd.Println(configTitle)
			cmd.Println()
			cmd.Println(strings.Join(configs, "\n"))
			cmd.Println()
			cmd.Println(separator)
			cmd.Println(functionTitle)
			cmd.Println()
			cmd.Println(strings.Join(functions, "\n"))
			cmd.Println()
			cmd.Println(separator)
			cmd.Println(envTitle)
			cmd.Println()
			cmd.Println(strings.Join(envs, "\n"))
			cmd.Println()
			cmd.Println(separator)
			return nil
		},
	}
}

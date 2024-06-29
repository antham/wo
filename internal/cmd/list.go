package cmd

import (
	"errors"
	"fmt"
	"strings"

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
			title := titleStyle.Render("Workspaces")
			var list []string
			for _, w := range workspaces {
				list = append(list, regularStyle.
					Render(fmt.Sprintf("* %s", w.Name)))
			}
			cmd.Println(title)
			cmd.Println()
			cmd.Println(separator)
			cmd.Println(strings.Join(list, "\n"))
			return nil
		},
	}
}

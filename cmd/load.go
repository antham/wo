package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var env string

func newLoadCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	loadCmd := &cobra.Command{
		Use:               "load workspace [environment]",
		Aliases:           []string{"l"},
		Short:             "Load a workspace",
		Args:              cobra.RangeArgs(1, 2),
		ValidArgsFunction: completionManager.Process,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 1:
				return workspaceManager.Load(args[0], "")
			case 2:
				return workspaceManager.Load(args[0], args[1])
			}
			return errors.New("too much arguments provided")
		},
	}
	return loadCmd
}

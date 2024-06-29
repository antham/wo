package cmd

import (
	"github.com/spf13/cobra"
)

func newConfigSetCmd(workspaceManager workspaceManager, completionManager completionManager) *cobra.Command {
	return &cobra.Command{
		Use:               "set workspace key value",
		Short:             "Set a configuration",
		Args:              cobra.ExactArgs(3),
		ValidArgsFunction: completionManager.Process,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := workspaceManager.SetConfig(args[0], map[string]string{args[1]: args[2]})
			if err != nil {
				return err
			}
			cmd.Printf("Config key '%s' edited on workspace '%s'\n", args[1], args[0])
			return nil
		},
	}
}

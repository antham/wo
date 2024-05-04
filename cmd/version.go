package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appVersion = "dev"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(appVersion)
		},
	}
}

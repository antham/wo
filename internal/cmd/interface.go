package cmd

import (
	"github.com/antham/wo/internal/workspace"
	"github.com/spf13/cobra"
)

type workspaceManager interface {
	BuildAliases(string) ([]string, error)
	Get(string) (workspace.Workspace, error)
	Create(string, string) error
	CreateEnv(string, string) error
	Edit(string) error
	EditEnv(string, string) error
	List() ([]workspace.Workspace, error)
	RunFunction(string, string, []string) error
	Remove(string) error
	SetConfig(string, map[string]string) error
	GetSupportedApps() []string
}

type completionManager interface {
	Process(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective)
}

package cmd

import (
	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
)

type workspaceManager interface {
	Get(string) (workspace.Workspace, error)
	Create(string, string) error
	CreateEnv(string, string) error
	Edit(string) error
	EditEnv(string, string) error
	Load(string, string) error
	List() ([]workspace.Workspace, error)
	RunFunction(string, string, []string) error
	Remove(string) error
	SetConfig(string, string, string) error
	Cd(string) error
}

type completionManager interface {
	Process(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective)
}

package cmd

import "github.com/antham/wo/workspace"

type workspaceManager interface {
	Get(string) (workspace.Workspace, error)
	Edit(string) error
	EditEnv(string, string) error
	Load(string, string) error
	List() ([]workspace.Workspace, error)
	RunFunction(string, string, []string) error
	Remove(string) error
}

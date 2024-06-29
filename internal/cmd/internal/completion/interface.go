package completion

import "github.com/antham/wo/internal/workspace"

type workspaceManager interface {
	List() ([]workspace.Workspace, error)
	Get(string) (workspace.Workspace, error)
	GetSupportedApps() []string
	GetConfigDir() string
}

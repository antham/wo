package completion

import "github.com/antham/wo/workspace"

type workspaceManager interface {
	List() ([]workspace.Workspace, error)
	Get(string) (workspace.Workspace, error)
}

package completion

import (
	"strings"

	"github.com/antham/wo/workspace"
)

type Completion struct {
	workspaceManager workspace.WorkspaceManager
}

func New(workspaceManager workspace.WorkspaceManager) Completion {
	return Completion{workspaceManager: workspaceManager}
}

func (c Completion) FindWorkspaces(toComplete string) ([]string, error) {
	workspaces, err := c.workspaceManager.List()
	if err != nil {
		return []string{}, err
	}
	ws := []string{}
	for _, w := range workspaces {
		if strings.HasPrefix(w.Name, toComplete) {
			ws = append(ws, w.Name)
		}
	}
	return ws, nil
}

func (c Completion) FindCommands(workspace string, toComplete string) ([]string, error) {
	w, err := c.workspaceManager.Get(workspace)
	if err != nil {
		return []string{}, err
	}
	cs := []string{}
	for _, c := range w.Commands {
		if strings.HasPrefix(c.Command, toComplete) {
			cs = append(cs, c.Command)
		}
	}
	return cs, nil
}

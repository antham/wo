package completion

import (
	"fmt"
	"strings"
)

type Completion struct {
	workspaceManager workspaceManager
}

func New(workspaceManager workspaceManager) Completion {
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

func (c Completion) FindFunctions(workspace string, toComplete string) ([]string, error) {
	w, err := c.workspaceManager.Get(workspace)
	if err != nil {
		return []string{}, err
	}
	fs := []string{}
	for _, f := range w.Functions {
		if strings.HasPrefix(f.Name, toComplete) {
			s := f.Name
			if f.Description != "" {
				s = fmt.Sprintf("%s\t%s", f.Name, f.Description)
			}
			fs = append(fs, s)
		}
	}
	return fs, nil
}

func (c Completion) FindEnvs(workspace string, toComplete string) ([]string, error) {
	w, err := c.workspaceManager.Get(workspace)
	if err != nil {
		return []string{}, err
	}
	envs := []string{}
	for _, env := range w.Envs {
		if strings.HasPrefix(env, toComplete) {
			envs = append(envs, env)
		}
	}
	return envs, nil
}

package completion

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type Decorator func(workspaceManager, string, ...string) ([]string, error)

type Completion struct {
	workspaceManager workspaceManager
	decorators       []Decorator
}

func New(workspaceManager workspaceManager, decorators []Decorator) Completion {
	return Completion{workspaceManager: workspaceManager, decorators: decorators}
}

func (c Completion) Process(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > len(c.decorators) {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	var matches []string
	var err error
	if len(args) == 0 {
		matches, err = c.decorators[0](c.workspaceManager, toComplete)
	} else {
		matches, err = c.decorators[len(args)](c.workspaceManager, toComplete, args[len(args)-1])
	}
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return matches, cobra.ShellCompDirectiveNoFileComp
}

func FindWorkspaces(workspaceManager workspaceManager, toComplete string, args ...string) ([]string, error) {
	workspaces, err := workspaceManager.List()
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

func FindFunctions(workspaceManager workspaceManager, toComplete string, args ...string) ([]string, error) {
	w, err := workspaceManager.Get(args[0])
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

func FindEnvs(workspaceManager workspaceManager, toComplete string, args ...string) ([]string, error) {
	w, err := workspaceManager.Get(args[0])
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

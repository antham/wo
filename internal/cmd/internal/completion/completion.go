package completion

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type Decorator func(workspaceManager, string, ...string) ([]string, cobra.ShellCompDirective, error)

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
	var shellDirective cobra.ShellCompDirective
	if len(args) == 0 {
		matches, shellDirective, err = c.decorators[0](c.workspaceManager, toComplete)
	} else {
		matches, shellDirective, err = c.decorators[len(args)](c.workspaceManager, toComplete, args[len(args)-1])
	}
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return matches, shellDirective
}

func FindWorkspaces(workspaceManager workspaceManager, toComplete string, args ...string) ([]string, cobra.ShellCompDirective, error) {
	workspaces, err := workspaceManager.List()
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp, err
	}
	ws := []string{}
	for _, w := range workspaces {
		if strings.HasPrefix(w.Name, toComplete) {
			ws = append(ws, w.Name)
		}
	}
	return ws, cobra.ShellCompDirectiveNoFileComp, nil
}

func FindFunctions(workspaceManager workspaceManager, toComplete string, args ...string) ([]string, cobra.ShellCompDirective, error) {
	w, err := workspaceManager.Get(args[0])
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp, err
	}
	fs := []string{}
	for _, f := range w.Functions.Functions {
		if strings.HasPrefix(f.Name, toComplete) {
			s := f.Name
			if f.Description != "" {
				s = fmt.Sprintf("%s\t%s", f.Name, f.Description)
			}
			fs = append(fs, s)
		}
	}
	return fs, cobra.ShellCompDirectiveNoFileComp, nil
}

func FindEnvs(workspaceManager workspaceManager, toComplete string, args ...string) ([]string, cobra.ShellCompDirective, error) {
	w, err := workspaceManager.Get(args[0])
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp, err
	}
	envs := []string{}
	for _, env := range w.Envs {
		if strings.HasPrefix(env.Name, toComplete) {
			envs = append(envs, env.Name)
		}
	}
	return envs, cobra.ShellCompDirectiveNoFileComp, nil
}

func FindDirs(workspaceManager workspaceManager, toComplete string, args ...string) ([]string, cobra.ShellCompDirective, error) {
	return []string{}, cobra.ShellCompDirectiveFilterDirs, nil
}

func NoOp(workspaceManager workspaceManager, toComplete string, args ...string) ([]string, cobra.ShellCompDirective, error) {
	return []string{}, cobra.ShellCompDirectiveNoFileComp, nil
}

var config = map[string]func(workspaceManager, string) ([]string, cobra.ShellCompDirective, error){
	"app": func(workspaceManager workspaceManager, toComplete string) ([]string, cobra.ShellCompDirective, error) {
		apps := []string{}
		for _, app := range workspaceManager.GetSupportedApps() {
			if strings.HasPrefix(app, toComplete) {
				apps = append(apps, app)
			}
		}
		return apps, cobra.ShellCompDirectiveNoFileComp, nil
	},
	"path": func(workspaceManager, string) ([]string, cobra.ShellCompDirective, error) {
		return []string{}, cobra.ShellCompDirectiveFilterDirs, nil
	},
}

func FindConfigKey(workspaceManager workspaceManager, toComplete string, args ...string) ([]string, cobra.ShellCompDirective, error) {
	keys := []string{}
	for key := range config {
		if strings.HasPrefix(key, toComplete) {
			keys = append(keys, key)
		}
	}
	return keys, cobra.ShellCompDirectiveNoFileComp, nil
}

func FindConfigValue(workspaceManager workspaceManager, toComplete string, args ...string) ([]string, cobra.ShellCompDirective, error) {
	f, ok := config[args[0]]
	if ok {
		return f(workspaceManager, toComplete)
	}
	return []string{}, cobra.ShellCompDirectiveNoFileComp, nil
}

func FindGlobalConfigKey(workspaceManager workspaceManager, toComplete string, args ...string) ([]string, cobra.ShellCompDirective, error) {
	keys := []string{}
	for _, key := range []string{"config-dir"} {
		if strings.HasPrefix(key, toComplete) {
			keys = append(keys, key)
		}
	}
	return keys, cobra.ShellCompDirectiveNoFileComp, nil
}

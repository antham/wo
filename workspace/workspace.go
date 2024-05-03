package workspace

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/antham/wo/shell"
)

const (
	configDir         = ".config/wo"
	envVariablePrefix = "WO"
)

const (
	bash = "bash"
	fish = "fish"
	sh   = "sh"
	zsh  = "zsh"
)

type Workspace struct {
	Name      string
	Functions []Function
	Envs      []string
}

type Function struct {
	Function    string
	Description string
}

type WorkspaceManager struct {
	editor      string
	shellBin    string
	shell       string
	homeDir     string
	execCommand func(...string) error
}

func NewWorkspaceManager() (WorkspaceManager, error) {
	editor := os.Getenv("EDITOR")
	visual := os.Getenv("VISUAL")
	s := WorkspaceManager{}
	switch {
	case editor != "":
		s.editor = editor
	case visual != "":
		s.editor = visual
	default:
		return WorkspaceManager{}, errors.New("no VISUAL or EDITOR environment variable found")
	}
	s.shellBin = os.Getenv("SHELL")
	s.shell = path.Base(s.shellBin)
	usr, err := user.Current()
	if err != nil {
		return WorkspaceManager{}, err
	}
	s.homeDir = usr.HomeDir
	s.execCommand = execCommand(s.shellBin)
	return s, s.createConfigFolder()
}

func (s WorkspaceManager) List() ([]Workspace, error) {
	workspaces := []Workspace{}
	err := filepath.Walk(s.getFunctionDir(), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		name := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
		workspace, err := s.Get(name)
		if err != nil {
			return err
		}
		workspaces = append(workspaces, workspace)
		return nil
	})
	if err != nil {
		return []Workspace{}, err
	}
	return workspaces, nil
}

func (s WorkspaceManager) Get(name string) (Workspace, error) {
	content, err := os.ReadFile(s.resolveFunctionFile(name))
	if os.IsNotExist(err) {
		return Workspace{}, errors.New("the workspace does not exist")
	}
	if err != nil {
		return Workspace{}, err
	}
	funcs, err := shell.Parse(s.shell, content)
	if err != nil {
		return Workspace{}, err
	}
	envs, err := s.listEnvs(name)
	if err != nil {
		return Workspace{}, err
	}
	commands := []Function{}
	for _, f := range funcs {
		commands = append(commands, Function{
			Function:    f.Name,
			Description: f.Description,
		})
	}
	slices.SortFunc(commands, func(a, b Function) int {
		return cmp.Compare(a.Function, b.Function)
	})
	return Workspace{
		Name:      name,
		Functions: commands,
		Envs:      envs,
	}, nil
}

func (s WorkspaceManager) Edit(name string) error {
	err := s.createWorkspaceEnvFolder(name)
	if err != nil {
		return err
	}
	err = s.createWorkspaceDefaultEnv(name)
	if err != nil {
		return err
	}
	return s.edit(s.resolveFunctionFile(name))
}

func (s WorkspaceManager) EditEnv(name string, env string) error {
	err := s.createWorkspaceEnvFolder(name)
	if err != nil {
		return err
	}
	err = s.createWorkspaceDefaultEnv(name)
	if err != nil {
		return err
	}
	return s.edit(s.resolveWorkspaceEnvFile(name, env))
}

func (s WorkspaceManager) Load(name string, env string) error {
	return s.execCommand(s.appendLoadStatement(name, env, []string{})...)
}

func (s WorkspaceManager) RunFunction(name string, env string, functionAndArgs []string) error {
	return s.execCommand(s.appendLoadStatement(name, env, functionAndArgs)...)
}

func (s WorkspaceManager) Remove(name string) error {
	_, err := s.Get(name)
	if err != nil {
		return err
	}
	return errors.Join(os.Remove(s.resolveFunctionFile(name)), os.RemoveAll(s.getWorkspaceEnvDir(name)))
}

func (s WorkspaceManager) appendLoadStatement(name string, env string, functionAndArgs []string) []string {
	data := []string{}
	data = append(data, s.createEnvVariableStatement(fmt.Sprintf("%s_NAME", envVariablePrefix), name))
	data = append(data, s.createEnvVariableStatement(fmt.Sprintf("%s_ENV", envVariablePrefix), s.resolveEnv(env)))
	envFile := s.resolveWorkspaceEnvFile(name, env)
	_, eerr := os.Stat(envFile)
	if eerr == nil {
		data = append(data, fmt.Sprintf("source %s", envFile))
	}
	data = append(data, fmt.Sprintf("source %s", s.resolveFunctionFile(name)))
	stmts := []string{}
	switch s.shell {
	case bash, sh, zsh:
		if len(functionAndArgs) > 0 {
			data = append(data, strings.Join(functionAndArgs, " "))
		}
		stmts = append(stmts, "-c", strings.Join(data, " && "))
	case fish:
		for _, d := range data {
			stmts = append(stmts, "-C", d)
		}
		if len(functionAndArgs) > 0 {
			stmts = append(stmts, "-c", strings.Join(functionAndArgs, " "))
		}
	}
	return stmts
}

func (s WorkspaceManager) edit(filepath string) error {
	return s.execCommand("-c", fmt.Sprintf("%s %s", s.editor, filepath))
}

func (s WorkspaceManager) resolveEnv(env string) string {
	if env == "" {
		return "default"
	}
	return env
}

func (s WorkspaceManager) listEnvs(name string) ([]string, error) {
	envs := []string{}
	dir := s.getWorkspaceEnvDir(name)
	file, err := os.Open(dir)
	if err != nil {
		return envs, err
	}
	fs, err := file.Readdir(-1)
	if err != nil {
		return envs, err
	}
	for _, f := range fs {
		envs = append(envs, strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())))
	}
	sort.Strings(envs)
	return envs, err
}

func (s WorkspaceManager) resolveFunctionFile(name string) string {
	return fmt.Sprintf("%s/%s.%s", s.getFunctionDir(), name, s.getExtension())
}

func (s WorkspaceManager) resolveWorkspaceEnvFile(name string, env string) string {
	return fmt.Sprintf("%s/%s.%s", s.getWorkspaceEnvDir(name), s.resolveEnv(env), s.getExtension())
}

func (s WorkspaceManager) getExtension() string {
	for _, shell := range []string{fish, bash, zsh, sh} {
		if strings.Contains(s.shellBin, shell) {
			return shell
		}
	}
	return ""
}

func (s WorkspaceManager) createConfigFolder() error {
	return errors.Join(os.MkdirAll(s.getConfigDir(), 0o777), os.MkdirAll(s.getFunctionDir(), 0o777), os.MkdirAll(s.getEnvDir(), 0o777))
}

func (s WorkspaceManager) createWorkspaceEnvFolder(name string) error {
	return os.MkdirAll(s.getWorkspaceEnvDir(name), 0o777)
}

func (s WorkspaceManager) createWorkspaceDefaultEnv(name string) error {
	_, err := os.OpenFile(s.resolveWorkspaceEnvFile(name, "default"), os.O_CREATE, 0o666)
	return err
}

func (s WorkspaceManager) getConfigDir() string {
	return fmt.Sprintf("%s/%s", s.homeDir, configDir)
}

func (s WorkspaceManager) getFunctionDir() string {
	return fmt.Sprintf("%s/functions/%s", s.getConfigDir(), s.shell)
}

func (s WorkspaceManager) getEnvDir() string {
	return fmt.Sprintf("%s/envs/%s", s.getConfigDir(), s.shell)
}

func (s WorkspaceManager) getWorkspaceEnvDir(name string) string {
	return fmt.Sprintf("%s/%s", s.getEnvDir(), name)
}

func (s WorkspaceManager) createEnvVariableStatement(name string, value string) string {
	switch s.shell {
	case bash, sh, zsh:
		return fmt.Sprintf("export %s=%s", name, value)
	case fish:
		return fmt.Sprintf("set -x -g %s %s", name, value)
	}
	return ""
}

func execCommand(shellBin string) func(args ...string) error {
	return func(args ...string) error {
		command := exec.Command(shellBin, args...)
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr
		return command.Run()
	}
}

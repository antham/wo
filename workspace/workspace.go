package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/antham/wo/parser"
)

const configDir = ".config/wo"

type ByName []Workspace

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type ByCommand []Command

func (a ByCommand) Len() int           { return len(a) }
func (a ByCommand) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCommand) Less(i, j int) bool { return a[i].Command < a[j].Command }

type Workspace struct {
	Name     string
	Commands []Command
}

type Command struct {
	Command     string
	Description string
}

type WorkspaceManager struct {
	editor   string
	shellBin string
	shell    string
	homeDir  string
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
		name := strings.Split(info.Name(), ".")[0]
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
	if err != nil {
		return Workspace{}, err
	}
	fs, err := parser.Parse(s.shell, content)
	if err != nil {
		return Workspace{}, err
	}
	commands := []Command{}
	for _, f := range fs {
		commands = append(commands, Command{
			Command:     f.Name,
			Description: f.Description,
		})
	}
	return Workspace{
		Name:     name,
		Commands: commands,
	}, nil
}

func (s WorkspaceManager) Edit(name string) error {
	return s.edit(s.resolveFunctionFile(name))
}

func (s WorkspaceManager) EditEnv(name string, env string) error {
	return s.edit(s.resolveEnvFile(name, env))
}

func (s WorkspaceManager) Load(name string, env string) error {
	return s.execCommand(s.appendLoadStatement(name, env)...)
}

func (s WorkspaceManager) RunFunction(name string, env string, functionAndArgs []string) error {
	return s.execCommand(s.appendLoadStatement(name, env, "-c", strings.Join(functionAndArgs, " "))...)
}

func (s WorkspaceManager) Remove(name string) error {
	return os.Remove(s.resolveFunctionFile(name))
}

func (s WorkspaceManager) appendLoadStatement(name string, env string, cmds ...string) []string {
	stmts := []string{}
	envFile := s.resolveEnvFile(name, env)
	_, eerr := os.Stat(envFile)
	if eerr == nil {
		stmts = append(stmts, "-C", fmt.Sprintf("source %s", envFile))
	}
	stmts = append(stmts, "-C", fmt.Sprintf("source %s", s.resolveFunctionFile(name)))
	stmts = append(stmts, cmds...)
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

func (s WorkspaceManager) resolveFunctionFile(name string) string {
	return fmt.Sprintf("%s/%s.%s", s.getFunctionDir(), name, s.getExtension())
}

func (s WorkspaceManager) resolveEnvFile(name string, env string) string {
	return fmt.Sprintf("%s/%s.%s.%s", s.getEnvDir(), name, s.resolveEnv(env), s.getExtension())
}

func (s WorkspaceManager) getExtension() string {
	for _, shell := range []string{"fish", "bash", "zsh", "sh"} {
		if strings.Contains(s.shellBin, shell) {
			return shell
		}
	}
	return ""
}

func (s WorkspaceManager) execCommand(args ...string) error {
	command := exec.Command(s.shellBin, args...)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	return command.Run()
}

func (s WorkspaceManager) createConfigFolder() error {
	return errors.Join(os.MkdirAll(s.getConfigDir(), 0777), os.MkdirAll(s.getFunctionDir(), 0777), os.MkdirAll(s.getEnvDir(), 0777))
}

func (s WorkspaceManager) getConfigDir() string {
	return fmt.Sprintf("%s/%s", s.homeDir, configDir)
}

func (s WorkspaceManager) getFunctionDir() string {
	return fmt.Sprintf("%s/functions/%s", s.getConfigDir(), s.shellBin)
}

func (s WorkspaceManager) getEnvDir() string {
	return fmt.Sprintf("%s/envs/%s", s.getConfigDir(), s.shellBin)
}

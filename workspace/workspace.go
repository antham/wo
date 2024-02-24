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

	"github.com/antham/wo/shell"
)

const configDir = ".config/wo"

type Workspace struct {
	Name     string
	Commands []Command
	Envs     []string
}

type Command struct {
	Command     string
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
	envs := []string{}
	err = filepath.Walk(s.getEnvDir(), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		s := strings.Split(info.Name(), ".")
		if s[0] == name {
			envs = append(envs, s[1])
		}
		return nil
	})
	if err != nil {
		return Workspace{}, err
	}
	commands := []Command{}
	for _, f := range funcs {
		commands = append(commands, Command{
			Command:     f.Name,
			Description: f.Description,
		})
	}
	return Workspace{
		Name:     name,
		Commands: commands,
		Envs:     envs,
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
	w, err := s.Get(name)
	if err != nil {
		return err
	}
	errs := []error{os.Remove(s.resolveFunctionFile(name))}
	for _, env := range w.Envs {
		errs = append(errs, os.Remove(s.resolveEnvFile(name, env)))
	}
	return errors.Join(errs...)
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

func (s WorkspaceManager) createConfigFolder() error {
	return errors.Join(os.MkdirAll(s.getConfigDir(), 0777), os.MkdirAll(s.getFunctionDir(), 0777), os.MkdirAll(s.getEnvDir(), 0777))
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

func execCommand(shellBin string) func(args ...string) error {
	return func(args ...string) error {
		command := exec.Command(shellBin, args...)
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr
		return command.Run()
	}
}

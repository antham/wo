package workspace

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
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
	editor string
	shell  string
}

func NewWorkspaceManager() (WorkspaceManager, error) {
	s := WorkspaceManager{}
	editor := os.Getenv("EDITOR")
	visual := os.Getenv("VISUAL")

	switch {
	case editor != "":
		s.editor = editor
	case visual != "":
		s.editor = visual
	default:
		return WorkspaceManager{}, errors.New("no VISUAL or EDITOR environment variable found")
	}
	s.shell = os.Getenv("SHELL")
	return s, s.CreateConfigFolder()
}

func (s WorkspaceManager) CreateConfigFolder() error {
	return errors.Join(os.MkdirAll(s.getConfigDir(), 0777), os.MkdirAll(s.getFunctionDir(), 0777), os.MkdirAll(s.getEnvDir(), 0777))
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
	readFile, err := os.Open(s.resolveFunctionFile(name))
	if err != nil {
		return Workspace{}, err
	}
	defer readFile.Close()
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	commands := []Command{}
	c := Command{}
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if regexp.MustCompile("^#").MatchString(line) {
			c.Description = strings.TrimSpace(strings.Trim(line, "#"))
		}
		if regexp.MustCompile(`^\s*function`).MatchString(line) {
			line = strings.TrimSpace(line)
			f := strings.Fields(line)
			c.Command = f[1]
			commands = append(commands, c)
		}
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
		if strings.Contains(s.shell, shell) {
			return shell
		}
	}
	return ""
}

func (s WorkspaceManager) execCommand(args ...string) error {
	command := exec.Command(s.shell, args...)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	return command.Run()
}

func (s WorkspaceManager) getConfigDir() string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	return homeDir + "/" + configDir
}

func (s WorkspaceManager) getFunctionDir() string {
	return s.getConfigDir() + "/functions"
}

func (s WorkspaceManager) getEnvDir() string {
	return s.getConfigDir() + "/envs"
}

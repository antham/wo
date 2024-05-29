package workspace

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"syscall"

	"github.com/antham/wo/shell"
	"github.com/spf13/viper"
)

const (
	defaultConfigDir  = ".config/wo"
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
	Name        string
	Description string
}

type WorkspaceManager struct {
	editor    string
	shellBin  string
	shell     string
	configDir string
	exec      Commander
}

func NewWorkspaceManager(options ...func(*WorkspaceManager)) (WorkspaceManager, error) {
	w := WorkspaceManager{}
	usr, err := user.Current()
	if err != nil {
		return WorkspaceManager{}, err
	}
	w.configDir = fmt.Sprintf("%s/%s", usr.HomeDir, defaultConfigDir)
	for _, o := range options {
		o(&w)
	}
	if w.editor == "" {
		return WorkspaceManager{}, errors.New("no editor defined")
	}
	w.exec = newCommand(w.shellBin)
	return w, w.createConfigFolder()
}

func WithEditor(editor string, visual string) func(*WorkspaceManager) {
	return func(w *WorkspaceManager) {
		switch {
		case editor != "":
			w.editor = editor
		case visual != "":
			w.editor = visual
		}
	}
}

func WithShellPath(shell string) func(*WorkspaceManager) {
	return func(w *WorkspaceManager) {
		w.shellBin = shell
		w.shell = path.Base(shell)
	}
}

func WithConfigPath(path string) func(*WorkspaceManager) {
	return func(w *WorkspaceManager) {
		w.configDir = path
	}
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
			Name:        f.Name,
			Description: f.Description,
		})
	}
	slices.SortFunc(commands, func(a, b Function) int {
		return cmp.Compare(a.Name, b.Name)
	})
	return Workspace{
		Name:      name,
		Functions: commands,
		Envs:      envs,
	}, nil
}

func (s WorkspaceManager) Create(name string, path string) error {
	err := s.createWorkspaceEnvFolder(name)
	if err != nil {
		return err
	}
	err = s.createFile(s.resolveFunctionFile(name))
	if err != nil {
		return err
	}
	err = s.createFile(s.resolveEnvFile(name, s.resolveEnv("")))
	if err != nil {
		return err
	}
	err = s.createFile(s.resolveConfigFile(name))
	if err != nil {
		return err
	}
	return s.SetConfig(name, "path", path)
}

func (s WorkspaceManager) CreateEnv(name string, env string) error {
	functionFile := s.resolveFunctionFile(name)
	_, err := os.Stat(functionFile)
	if err != nil {
		return fmt.Errorf(`check the workspace "%s" exists, create it first`, name)
	}
	return s.createFile(s.resolveEnvFile(name, env))
}

func (s WorkspaceManager) Edit(name string) error {
	functionFile := s.resolveFunctionFile(name)
	_, err := os.Stat(functionFile)
	if err != nil {
		return fmt.Errorf(`check the workspace "%s" exists`, name)
	}
	return s.editFile(functionFile)
}

func (s WorkspaceManager) EditEnv(name string, env string) error {
	functionFile := s.resolveFunctionFile(name)
	_, err := os.Stat(functionFile)
	if err != nil {
		return fmt.Errorf(`check the workspace "%s" exists`, name)
	}
	envFile := s.resolveEnvFile(name, env)
	_, err = os.Stat(envFile)
	if err != nil {
		return fmt.Errorf(`check the environment "%s" exists`, env)
	}
	return s.editFile(envFile)
}

func (s WorkspaceManager) RunFunction(name string, env string, functionAndArgs []string) error {
	p, err := s.GetConfig(name, "path")
	if err != nil {
		return err
	}
	path := ""
	if p != nil {
		path = p.(string)
	}
	return s.exec.command(path, s.appendLoadStatement(name, env, functionAndArgs)...)
}

func (s WorkspaceManager) Remove(name string) error {
	_, err := s.Get(name)
	if err != nil {
		return err
	}
	return errors.Join(os.Remove(s.resolveFunctionFile(name)), os.RemoveAll(s.getWorkspaceEnvDir(name)))
}

func (s WorkspaceManager) SetConfig(name string, key string, value string) error {
	v := s.getViper(name)
	v.Set(key, value)
	return v.WriteConfig()
}

func (s WorkspaceManager) GetConfig(name string, key string) (any, error) {
	v := s.getViper(name)
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return v.Get(key), nil
}

func (s WorkspaceManager) Cd(name string) error {
	p, err := s.GetConfig(name, "path")
	if err != nil {
		return err
	}
	err = os.Chdir(p.(string))
	if err != nil {
		return err
	}
	return syscall.Exec(s.shellBin, []string{""}, os.Environ())
}

func (s WorkspaceManager) appendLoadStatement(name string, env string, functionAndArgs []string) []string {
	data := []string{}
	data = append(data, s.createEnvVariableStatement(fmt.Sprintf("%s_NAME", envVariablePrefix), name))
	data = append(data, s.createEnvVariableStatement(fmt.Sprintf("%s_ENV", envVariablePrefix), s.resolveEnv(env)))
	envFile := s.resolveEnvFile(name, env)
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

func (s WorkspaceManager) editFile(filepath string) error {
	return s.exec.command("", "-c", fmt.Sprintf("%s %s", s.editor, filepath))
}

func (s WorkspaceManager) createFile(filepath string) error {
	_, err := os.OpenFile(filepath, os.O_CREATE, 0o666)
	return err
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

func (s WorkspaceManager) resolveEnvFile(name string, env string) string {
	return fmt.Sprintf("%s/%s.%s", s.getWorkspaceEnvDir(name), s.resolveEnv(env), s.getExtension())
}

func (s WorkspaceManager) resolveConfigFile(name string) string {
	return fmt.Sprintf("%s/%s.toml", s.getConfigDir(), name)
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
	return errors.Join(
		os.MkdirAll(s.configDir, 0o777),
		os.MkdirAll(s.getFunctionDir(), 0o777),
		os.MkdirAll(s.getEnvDir(), 0o777),
		os.MkdirAll(s.getConfigDir(), 0o777),
	)
}

func (s WorkspaceManager) createWorkspaceEnvFolder(name string) error {
	return os.MkdirAll(s.getWorkspaceEnvDir(name), 0o777)
}

func (s WorkspaceManager) getFunctionDir() string {
	return fmt.Sprintf("%s/functions", s.configDir)
}

func (s WorkspaceManager) getEnvDir() string {
	return fmt.Sprintf("%s/envs", s.configDir)
}

func (s WorkspaceManager) getConfigDir() string {
	return fmt.Sprintf("%s/configs", s.configDir)
}

func (s WorkspaceManager) getWorkspaceEnvDir(name string) string {
	return fmt.Sprintf("%s/%s", s.getEnvDir(), name)
}

func (s WorkspaceManager) getViper(name string) *viper.Viper {
	v := viper.New()
	v.AddConfigPath(fmt.Sprintf("%s/", s.getConfigDir()))
	v.SetConfigName(name)
	v.SetConfigType("toml")
	return v
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

type command struct {
	shellBin string
}

func newCommand(shellBin string) *command {
	return &command{
		shellBin: shellBin,
	}
}

func (c *command) command(path string, args ...string) error {
	command := exec.Command(c.shellBin, args...)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	command.Dir = path
	slog.With(slog.String("command", command.String())).With(slog.String("path", command.Dir)).Debug("command to run")
	return command.Run()
}

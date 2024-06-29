package workspace

import (
	"cmp"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/antham/wo/internal/shell"
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
	Functions Functions
	Envs      []Env
	Config    map[string]string
	dir       string
}

type Functions struct {
	file      string
	Functions []Function
}

type Function struct {
	Name        string
	Description string
}

type Env struct {
	Name string
	file string
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

func (s WorkspaceManager) BuildAliases(prefix string) ([]string, error) {
	workspaces, err := s.List()
	if err != nil {
		return []string{}, err
	}
	aliases := []string{}
	for _, w := range workspaces {
		aliases = append(aliases, fmt.Sprintf(`alias %s%s="cd %s"`, prefix, w.Name, w.Config["path"]))
	}
	return aliases, nil
}

func (s WorkspaceManager) List() ([]Workspace, error) {
	workspaces := []Workspace{}
	entries, err := os.ReadDir(s.configDir)
	if err != nil {
		return workspaces, err
	}
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		workspace, err := s.getWorkspace(e.Name())
		if err != nil {
			return workspaces, err
		}
		workspaces = append(workspaces, workspace)
	}
	return workspaces, nil
}

func (s WorkspaceManager) Get(name string) (Workspace, error) {
	return s.getWorkspace(name)
}

func (s WorkspaceManager) Create(name string, path string) error {
	err := s.createWorkspaceFolder(name)
	if err != nil {
		return err
	}
	err = s.createFile(s.resolveFunctionFile(name))
	if err != nil {
		return err
	}
	err = s.createFile(s.resolveEnvFile(name, "default"))
	if err != nil {
		return err
	}
	err = s.createFile(s.resolveConfigFile(name))
	if err != nil {
		return err
	}
	return s.SetConfig(
		name,
		map[string]string{
			"app":  s.shell,
			"path": path,
		},
	)
}

func (s WorkspaceManager) CreateEnv(name string, env string) error {
	_, err := s.getWorkspace(name)
	if err != nil {
		return err
	}
	return s.createFile(s.resolveEnvFile(name, env))
}

func (s WorkspaceManager) Edit(name string) error {
	w, err := s.getWorkspace(name)
	if err != nil {
		return err
	}
	return s.editFile(w.Functions.file)
}

func (s WorkspaceManager) EditEnv(name string, env string) error {
	w, err := s.getWorkspace(name)
	if err != nil {
		return err
	}
	index := slices.IndexFunc(w.Envs, func(e Env) bool {
		return e.Name == env
	})
	if index == -1 {
		return fmt.Errorf("the env `%s` does not exist", env)
	}
	return s.editFile(w.Envs[index].file)
}

func (s WorkspaceManager) RunFunction(name string, env string, functionAndArgs []string) error {
	w, err := s.getWorkspace(name)
	if err != nil {
		return err
	}
	if !slices.ContainsFunc(w.Envs, func(e Env) bool {
		return e.Name == env
	}) {
		return fmt.Errorf("the env `%s` does not exist", env)
	}
	if !slices.ContainsFunc(w.Functions.Functions, func(f Function) bool {
		return f.Name == functionAndArgs[0]
	}) {
		return fmt.Errorf("the function `%s` does not exist", functionAndArgs[0])
	}
	return s.exec.command(w.Config["path"], s.appendLoadStatement(name, env, functionAndArgs)...)
}

func (s WorkspaceManager) Remove(name string) error {
	w, err := s.Get(name)
	if err != nil {
		return err
	}
	return os.RemoveAll(w.dir)
}

func (s WorkspaceManager) SetConfig(name string, kv map[string]string) error {
	v := s.getViper(name)
	err := v.ReadInConfig()
	if err != nil {
		return err
	}
	for key, value := range kv {
		if !slices.Contains([]string{"path", "app"}, key) {
			return fmt.Errorf(`"%s" is not a valid config key`, key)
		}
		if key == "path" {
			_, err := os.Stat(value)
			if os.IsNotExist(err) {
				return fmt.Errorf(`path "%s" does not exist`, value)
			}
		}
		if key == "app" {
			if !slices.Contains(s.GetSupportedApps(), value) {
				return fmt.Errorf(`app "%s" is not supported`, value)
			}
		}
		v.Set(key, value)
	}
	return v.WriteConfig()
}

func (s WorkspaceManager) GetConfig(name string, key string) (string, error) {
	v := s.getViper(name)
	err := v.ReadInConfig()
	if err != nil {
		return "", err
	}
	return v.GetString(key), nil
}

func (s WorkspaceManager) GetSupportedApps() []string {
	return []string{fish, bash, zsh, sh}
}

func (s WorkspaceManager) GetConfigDir() string {
	return s.configDir
}

func (s WorkspaceManager) appendLoadStatement(name string, env string, functionAndArgs []string) []string {
	data := []string{}
	data = append(data, s.createEnvVariableStatement(fmt.Sprintf("%s_NAME", envVariablePrefix), name))
	data = append(data, s.createEnvVariableStatement(fmt.Sprintf("%s_ENV", envVariablePrefix), env))
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

func (s WorkspaceManager) listEnvs(name string) ([]Env, error) {
	envs := []Env{}
	dir := s.getWorkspaceEnvsDir(name)
	file, err := os.Open(dir)
	if err != nil {
		return []Env{}, err
	}
	fs, err := file.Readdir(-1)
	if err != nil {
		return []Env{}, err
	}
	for _, f := range fs {
		env := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		envs = append(envs, Env{Name: env, file: s.resolveEnvFile(name, env)})
	}
	sort.Slice(envs, func(i, j int) bool {
		return envs[i].Name < envs[j].Name
	})
	return envs, err
}

func (s WorkspaceManager) resolveFunctionFile(name string) string {
	return fmt.Sprintf("%s/functions.%s", s.getWorkspaceFunctionsDir(name), s.getExtension())
}

func (s WorkspaceManager) resolveEnvFile(name string, env string) string {
	return fmt.Sprintf("%s/%s.%s", s.getWorkspaceEnvsDir(name), env, s.getExtension())
}

func (s WorkspaceManager) resolveConfigFile(name string) string {
	return fmt.Sprintf("%s/config.toml", s.getWorkspaceDir(name))
}

func (s WorkspaceManager) getExtension() string {
	for _, shell := range s.GetSupportedApps() {
		if strings.Contains(s.shellBin, shell) {
			return shell
		}
	}
	return ""
}

func (s WorkspaceManager) createConfigFolder() error {
	return os.MkdirAll(s.configDir, 0o777)
}

func (s WorkspaceManager) createWorkspaceFolder(name string) error {
	return errors.Join(
		os.MkdirAll(s.getWorkspaceFunctionsDir(name), 0o777),
		os.MkdirAll(s.getWorkspaceEnvsDir(name), 0o777),
	)
}

func (s WorkspaceManager) getWorkspaceDir(name string) string {
	return fmt.Sprintf("%s/%s", s.configDir, name)
}

func (s WorkspaceManager) getWorkspaceFunctionsDir(name string) string {
	return fmt.Sprintf("%s/functions", s.getWorkspaceDir(name))
}

func (s WorkspaceManager) getWorkspaceEnvsDir(name string) string {
	return fmt.Sprintf("%s/envs", s.getWorkspaceDir(name))
}

func (s WorkspaceManager) getViper(name string) *viper.Viper {
	v := viper.New()
	v.AddConfigPath(fmt.Sprintf("%s/", s.getWorkspaceDir(name)))
	v.SetConfigName("config")
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

func (s WorkspaceManager) getWorkspace(name string) (Workspace, error) {
	_, err := os.Stat(s.getWorkspaceDir(name))
	if os.IsNotExist(err) {
		return Workspace{}, errors.New("the workspace does not exist")
	}
	app, err := s.GetConfig(name, "app")
	if err != nil {
		return Workspace{}, errors.New("the config file of the workspace is corrupted")
	}
	path, err := s.GetConfig(name, "path")
	if err != nil {
		return Workspace{}, errors.New("the config file of the workspace is corrupted")
	}
	if app != s.shell {
		return Workspace{}, fmt.Errorf(`the "%s" app si not supported for this workspace, it works with "%s"`, app, s.shell)
	}
	content, err := os.ReadFile(s.resolveFunctionFile(name))
	if os.IsNotExist(err) {
		return Workspace{}, errors.New("the workspace does not exist")
	}
	if err != nil {
		return Workspace{}, err
	}
	funcs := shell.Parse(s.shell, content)
	envs, err := s.listEnvs(name)
	if err != nil {
		return Workspace{}, err
	}
	functions := []Function{}
	for _, f := range funcs {
		functions = append(
			functions, Function{
				Name:        f.Name,
				Description: f.Description,
			},
		)
	}
	slices.SortFunc(functions, func(a, b Function) int {
		return cmp.Compare(a.Name, b.Name)
	})
	return Workspace{
		Name: name,
		Functions: Functions{
			file:      s.resolveFunctionFile(name),
			Functions: functions,
		},
		Envs: envs,
		Config: map[string]string{
			"path": path,
			"app":  app,
		},
		dir: s.getWorkspaceDir(name),
	}, nil
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

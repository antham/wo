package workspace

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	configPath  string
	projectPath string
)

func getConfigPath(t *testing.T) string {
	if configPath == "" {
		path, err := os.MkdirTemp("/tmp", "wo")
		assert.NoError(t, err)
		configPath = path
	}
	return configPath
}

func getProjectPath(t *testing.T) string {
	if projectPath == "" {
		path, err := os.MkdirTemp("/tmp", "project")
		assert.NoError(t, err)
		projectPath = path
	}
	return projectPath
}

func TestNewWorkspaceManager(t *testing.T) {
	type scenario struct {
		name  string
		setup func() (string, string, string)
		test  func(*testing.T, WorkspaceManager, error)
	}
	scenarios := []scenario{
		{
			"No variable editor, nor visual defined",
			func() (string, string, string) {
				return "", "", ""
			},
			func(t *testing.T, w WorkspaceManager, err error) {
				assert.EqualError(t, err, "no editor defined")
			},
		},
		{
			"visual defined",
			func() (string, string, string) {
				return "vim", "", "/usr/bin/bash"
			},
			func(t *testing.T, w WorkspaceManager, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "vim", w.editor)
				assert.Equal(t, "/usr/bin/bash", w.shellBin)
				assert.Equal(t, "bash", w.shell)
				assert.DirExists(t, w.configDir)
				assert.DirExists(t, w.getFunctionDir())
				assert.DirExists(t, w.getEnvDir())
			},
		},
		{
			"editor defined",
			func() (string, string, string) {
				return "", "emacs", "/bin/zsh"
			},
			func(t *testing.T, w WorkspaceManager, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "emacs", w.editor)
				assert.Equal(t, "/bin/zsh", w.shellBin)
				assert.Equal(t, "zsh", w.shell)
				assert.DirExists(t, w.configDir)
				assert.DirExists(t, w.getFunctionDir())
				assert.DirExists(t, w.getEnvDir())
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			visual, editor, shell := s.setup()
			manager, err := NewWorkspaceManager(WithEditor(editor, visual), WithShellPath(shell), WithConfigPath(getConfigPath(t)))
			s.test(t, manager, err)
		})
	}
}

func TestList(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T, WorkspaceManager)
		test  func(*testing.T, []Workspace, error)
	}
	scenarios := []scenario{
		{
			"No workspace defined",
			func(t *testing.T, w WorkspaceManager) {
			},
			func(t *testing.T, ws []Workspace, err error) {
				assert.NoError(t, err)
				assert.Empty(t, ws)
			},
		},
		{
			"Get all workspaces ordered alphabetically",
			func(t *testing.T, w WorkspaceManager) {
				assert.NoError(t, w.Create("api", getProjectPath(t)))
				assert.NoError(t, w.CreateEnv("api", "dev"))
				assert.NoError(t, w.Create("db", getProjectPath(t)))
				assert.NoError(t, w.CreateEnv("db", "staging"))
				assert.NoError(t, w.Create("front", getProjectPath(t)))
				assert.NoError(t, w.CreateEnv("front", "prod"))
			},
			func(t *testing.T, ws []Workspace, err error) {
				assert.NoError(t, err)
				assert.Len(t, ws, 3)
				assert.Equal(t, []Workspace{
					{Name: "api", Functions: []Function{}, Envs: []string{"default", "dev"}},
					{Name: "db", Functions: []Function{}, Envs: []string{"default", "staging"}},
					{Name: "front", Functions: []Function{}, Envs: []string{"default", "prod"}},
				}, ws)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(getConfigPath(t)))
			assert.NoError(t, err)
			s.setup(t, w)
			workspaces, err := w.List()
			s.test(t, workspaces, err)
		})
	}
}

func TestGet(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T, WorkspaceManager)
		test  func(*testing.T, Workspace, error)
	}
	scenarios := []scenario{
		{
			"Workspace does not exist",
			func(t *testing.T, w WorkspaceManager) {
			},
			func(t *testing.T, w Workspace, err error) {
				assert.EqualError(t, err, "the workspace does not exist")
			},
		},
		{
			"Get all workspace",
			func(t *testing.T, w WorkspaceManager) {
				err := w.Create("front", getProjectPath(t))
				assert.NoError(t, err)
				err = w.CreateEnv("front", "prod")
				assert.NoError(t, err)

				functionPath := getConfigPath(t) + "/functions/front.bash"
				assert.NoError(t, os.WriteFile(functionPath, []byte(`
# A function 1
test_func1() {

}

# A function 2
test_func2() {

}
`), 0o777))
			},
			func(t *testing.T, w Workspace, err error) {
				assert.NoError(t, err)
				assert.Len(t, w.Functions, 2)
				assert.Equal(t, Workspace{
					Name: "front",
					Functions: []Function{
						{
							Name:        "test_func1",
							Description: "A function 1",
						},
						{
							Name:        "test_func2",
							Description: "A function 2",
						},
					},
					Envs: []string{"default", "prod"},
				}, w)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(getConfigPath(t)))
			assert.NoError(t, err)
			s.setup(t, w)
			workspace, err := w.Get("front")
			s.test(t, workspace, err)
		})
	}
}

func TestCreate(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T, WorkspaceManager)
		test  func(*testing.T, error)
	}
	scenarios := []scenario{
		{
			"Create workspace",
			func(t *testing.T, w WorkspaceManager) {
			},
			func(t *testing.T, err error) {
				assert.NoError(t, err)
				path := getConfigPath(t)
				envFile, err := os.Stat(path + "/envs/test/default.bash")
				assert.NoError(t, err)
				assert.Equal(t, "default.bash", envFile.Name())
				functionFile, err := os.Stat(path + "/functions/test.bash")
				assert.NoError(t, err)
				assert.Equal(t, "test.bash", functionFile.Name())
				configFile, err := os.Stat(path + "/configs/test.toml")
				assert.NoError(t, err)
				assert.Equal(t, "test.toml", configFile.Name())
				b, err := os.ReadFile(path + "/configs/test.toml")
				assert.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("path = '%s'\n", getProjectPath(t)), string(b))
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(getConfigPath(t)))
			assert.NoError(t, err)
			s.setup(t, w)
			assert.NoError(t, err)
			s.test(t, w.Create("test", getProjectPath(t)))
		})
	}
}

func TestCreateEnv(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T, WorkspaceManager)
		test  func(*testing.T, error)
	}
	scenarios := []scenario{
		{
			"Create workspace env",
			func(t *testing.T, w WorkspaceManager) {
			},
			func(t *testing.T, err error) {
				assert.NoError(t, err)
				path := getConfigPath(t)
				defaultEnvFile, err := os.Stat(path + "/envs/test/default.bash")
				assert.NoError(t, err)
				assert.Equal(t, "default.bash", defaultEnvFile.Name())
				prodEnvFile, err := os.Stat(path + "/envs/test/prod.bash")
				assert.NoError(t, err)
				assert.Equal(t, "prod.bash", prodEnvFile.Name())
				functionFile, err := os.Stat(path + "/functions/test.bash")
				assert.NoError(t, err)
				assert.Equal(t, "test.bash", functionFile.Name())
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(getConfigPath(t)))
			assert.NoError(t, err)
			s.setup(t, w)
			assert.NoError(t, err)
			assert.NoError(t, w.Create("test", getProjectPath(t)))
			s.test(t, w.CreateEnv("test", "prod"))
		})
	}
}

func TestEdit(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T, WorkspaceManager, *MockCommander)
		test  func(*testing.T)
	}
	scenarios := []scenario{
		{
			"Edit workspace",
			func(t *testing.T, w WorkspaceManager, exec *MockCommander) {
				exec.On("command", "", "-c", fmt.Sprintf("emacs %s/functions/test.bash", getConfigPath(t))).Return(nil)
				err := w.Create("test", getProjectPath(t))
				assert.NoError(t, err)
			},
			func(t *testing.T) {
				f, err := os.Stat(fmt.Sprintf("%s/envs/test/default.bash", getConfigPath(t)))
				assert.NoError(t, err)
				assert.Equal(t, "default.bash", f.Name())
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(getConfigPath(t)))
			assert.NoError(t, err)
			exec := NewMockCommander(t)
			w.exec = exec
			s.setup(t, w, exec)
			assert.NoError(t, w.Edit("test"))
			s.test(t)
		})
	}
}

func TestEditEnv(t *testing.T) {
	type scenario struct {
		name  string
		env   string
		setup func(*testing.T, *MockCommander)
	}
	scenarios := []scenario{
		{
			"Edit default workspace",
			"",
			func(t *testing.T, exec *MockCommander) {
				exec.On("command", "", "-c", fmt.Sprintf("emacs %s/envs/test/default.bash", getConfigPath(t))).Return(nil)
			},
		},
		{
			"Edit prod workspace",
			"prod",
			func(t *testing.T, exec *MockCommander) {
				exec.On("command", "", "-c", fmt.Sprintf("emacs %s/envs/test/prod.bash", getConfigPath(t))).Return(nil)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(getConfigPath(t)))
			assert.NoError(t, err)
			err = w.Create("test", getProjectPath(t))
			assert.NoError(t, err)
			err = w.CreateEnv("test", "prod")
			assert.NoError(t, err)
			exec := NewMockCommander(t)
			w.exec = exec
			s.setup(t, exec)
			assert.NoError(t, w.EditEnv("test", s.env))
		})
	}
}

func TestRunFunction(t *testing.T) {
	type scenario struct {
		name            string
		functionAndArgs []string
		env             string
		shell           string
		setup           func(*testing.T, *MockCommander)
	}
	scenarios := []scenario{
		{
			"Run a function with a bash shell",
			[]string{"run-db"},
			"",
			"/bin/bash",
			func(t *testing.T, exec *MockCommander) {
				exec.On("command", getProjectPath(t), "-c", fmt.Sprintf("export WO_NAME=test && export WO_ENV=default && source %s/envs/test/default.bash && source %s/functions/test.bash && run-db", getConfigPath(t), getConfigPath(t))).Return(nil)
			},
		},
		{
			"Run a function with a fish shell",
			[]string{"run-db"},
			"",
			"/bin/fish",
			func(t *testing.T, exec *MockCommander) {
				exec.On("command", getProjectPath(t), "-C", "set -x -g WO_NAME test", "-C", "set -x -g WO_ENV default", "-C", fmt.Sprintf("source %s/envs/test/default.fish", getConfigPath(t)), "-C", fmt.Sprintf("source %s/functions/test.fish", getConfigPath(t)), "-c", "run-db").Return(nil)
			},
		},
		{
			"Run a function with an env defined and a bash shell",
			[]string{"run-db", "watch"},
			"prod",
			"/bin/bash",
			func(t *testing.T, exec *MockCommander) {
				exec.On("command", getProjectPath(t), "-c", fmt.Sprintf("export WO_NAME=test && export WO_ENV=prod && source %s/envs/test/prod.bash && source %s/functions/test.bash && run-db watch", getConfigPath(t), getConfigPath(t))).Return(nil)
			},
		},
		{
			"Run a function with an env defined and a fish shell",
			[]string{"run-db", "watch"},
			"prod",
			"/bin/fish",
			func(t *testing.T, exec *MockCommander) {
				exec.On("command", getProjectPath(t), "-C", "set -x -g WO_NAME test", "-C", "set -x -g WO_ENV prod", "-C", fmt.Sprintf("source %s/envs/test/prod.fish", getConfigPath(t)), "-C", fmt.Sprintf("source %s/functions/test.fish", getConfigPath(t)), "-c", "run-db watch").Return(nil)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath(s.shell), WithConfigPath(getConfigPath(t)))
			assert.NoError(t, err)
			err = w.Create("test", getProjectPath(t))
			assert.NoError(t, err)
			err = w.CreateEnv("test", "prod")
			assert.NoError(t, err)
			exec := NewMockCommander(t)
			w.exec = exec
			s.setup(t, exec)
			assert.NoError(t, w.RunFunction("test", s.env, s.functionAndArgs))
		})
	}
}

func TestRemove(t *testing.T) {
	type scenario struct {
		name      string
		workspace string
		test      func(*testing.T, error)
	}
	scenarios := []scenario{
		{
			"Remove an unexisting workspace",
			"whatever",
			func(t *testing.T, e error) {
				assert.Error(t, e)
			},
		},
		{
			"Remove a workspace",
			"test",
			func(t *testing.T, e error) {
				assert.NoError(t, e)
				path := getConfigPath(t)
				_, err := os.Stat(path + "/envs/test/prod.bash")
				assert.True(t, os.IsNotExist(err))
				_, err = os.Stat(path + "/envs/test/dev.bash")
				assert.True(t, os.IsNotExist(err))
				_, err = os.Stat(path + "/functions/test.bash")
				assert.True(t, os.IsNotExist(err))

				_, err = os.Stat(path + "/functions/front.bash")
				assert.NoError(t, err)
				_, err = os.Stat(path + "/envs/front/dev.bash")
				assert.NoError(t, err)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(getConfigPath(t)))
			assert.NoError(t, err)
			err = w.Create("test", getProjectPath(t))
			assert.NoError(t, err)
			err = w.CreateEnv("test", "prod")
			assert.NoError(t, err)
			err = w.CreateEnv("test", "dev")
			assert.NoError(t, err)
			err = w.Create("front", getProjectPath(t))
			assert.NoError(t, err)
			err = w.CreateEnv("front", "dev")
			assert.NoError(t, err)
			s.test(t, w.Remove(s.workspace))
		})
	}
}

func TestSetConfig(t *testing.T) {
	type scenario struct {
		name      string
		workspace string
		key       string
		value     string
		test      func(*testing.T, error)
	}
	scenarios := []scenario{
		{
			"Set a value in a workspace config",
			"test",
			"path",
			"/home/user/project",
			func(t *testing.T, err error) {
				assert.NoError(t, err)
				b, err := os.ReadFile(fmt.Sprintf("%s/%s", getConfigPath(t), "configs/test.toml"))
				assert.NoError(t, err)
				assert.Equal(t, []byte("path = '/home/user/project'\n"), b)
			},
		},
		{
			"Set a value in an unexisting workspace",
			"whatever",
			"path",
			"/home/user/project",
			func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(getConfigPath(t)))
			assert.NoError(t, err)
			err = w.Create("test", getProjectPath(t))
			assert.NoError(t, err)
			s.test(t, w.SetConfig(s.workspace, s.key, s.value))
		})
	}
}

func TestBuildAliases(t *testing.T) {
	type scenario struct {
		name   string
		prefix string
		test   func(*testing.T, []string, error)
	}
	scenarios := []scenario{
		{
			"Build aliases for all workspace",
			"",
			func(t *testing.T, aliases []string, e error) {
				assert.NoError(t, e)
				assert.Equal(t, []string{
					fmt.Sprintf(`alias c_front="cd %s/front"`, getProjectPath(t)),
					fmt.Sprintf(`alias c_test="cd %s/test"`, getProjectPath(t)),
				}, aliases)
			},
		},
		{
			"Build aliases with a different prefix",
			"xx",
			func(t *testing.T, aliases []string, e error) {
				assert.NoError(t, e)
				assert.Equal(t, []string{
					fmt.Sprintf(`alias xxfront="cd %s/front"`, getProjectPath(t)),
					fmt.Sprintf(`alias xxtest="cd %s/test"`, getProjectPath(t)),
				}, aliases)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(getConfigPath(t)))
			assert.NoError(t, err)
			err = w.Create("test", fmt.Sprintf("%s/test", getProjectPath(t)))
			assert.NoError(t, err)
			err = w.Create("front", fmt.Sprintf("%s/front", getProjectPath(t)))
			assert.NoError(t, err)
			aliases, err := w.BuildAliases(s.prefix)
			s.test(t, aliases, err)
		})
	}
}
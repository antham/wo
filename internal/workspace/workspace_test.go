package workspace

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type config struct {
	path string
}

func (c *config) getPath(t *testing.T) string {
	if c.path == "" {
		path, err := os.MkdirTemp("/tmp", "wo")
		assert.NoError(t, err)
		c.path = path
	}
	return c.path
}

type project struct {
	path string
}

func (p *project) getPath(t *testing.T) string {
	if p.path == "" {
		path, err := os.MkdirTemp("/tmp", "project")
		assert.NoError(t, err)
		p.path = path
	}
	return p.path
}

func TestNewWorkspaceManager(t *testing.T) {
	config := &config{}
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
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(config.getPath(t))
			visual, editor, shell := s.setup()
			manager, err := NewWorkspaceManager(WithEditor(editor, visual), WithShellPath(shell), WithConfigPath(config.getPath(t)))
			s.test(t, manager, err)
		})
	}
}

func TestList(t *testing.T) {
	config := &config{}
	project := &project{}
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
				assert.NoError(t, w.Create("api", project.getPath(t)))
				assert.NoError(t, w.CreateEnv("api", "dev"))
				assert.NoError(t, w.Create("db", project.getPath(t)))
				assert.NoError(t, w.CreateEnv("db", "staging"))
				assert.NoError(t, w.Create("front", project.getPath(t)))
				assert.NoError(t, w.CreateEnv("front", "prod"))
			},
			func(t *testing.T, ws []Workspace, err error) {
				assert.NoError(t, err)
				assert.Len(t, ws, 3)
				assert.Equal(t, []Workspace{
					{Name: "api", Functions: Functions{file: fmt.Sprintf("%s/api/functions/functions.bash", config.getPath(t)), Functions: []Function{}}, Envs: []Env{{Name: "default", file: fmt.Sprintf("%s/api/envs/default.bash", config.getPath(t))}, {Name: "dev", file: fmt.Sprintf("%s/api/envs/dev.bash", config.getPath(t))}}, Config: map[string]string{"app": "bash", "path": project.getPath(t)}, dir: fmt.Sprintf("%s/api", config.getPath(t))},
					{Name: "db", Functions: Functions{file: fmt.Sprintf("%s/db/functions/functions.bash", config.getPath(t)), Functions: []Function{}}, Envs: []Env{{Name: "default", file: fmt.Sprintf("%s/db/envs/default.bash", config.getPath(t))}, {Name: "staging", file: fmt.Sprintf("%s/db/envs/staging.bash", config.getPath(t))}}, Config: map[string]string{"app": "bash", "path": project.getPath(t)}, dir: fmt.Sprintf("%s/db", config.getPath(t))},
					{Name: "front", Functions: Functions{file: fmt.Sprintf("%s/front/functions/functions.bash", config.getPath(t)), Functions: []Function{}}, Envs: []Env{{Name: "default", file: fmt.Sprintf("%s/front/envs/default.bash", config.getPath(t))}, {Name: "prod", file: fmt.Sprintf("%s/front/envs/prod.bash", config.getPath(t))}}, Config: map[string]string{"app": "bash", "path": project.getPath(t)}, dir: fmt.Sprintf("%s/front", config.getPath(t))},
				}, ws)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(config.getPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(config.getPath(t)))
			assert.NoError(t, err)
			s.setup(t, w)
			workspaces, err := w.List()
			s.test(t, workspaces, err)
		})
	}
}

func TestGet(t *testing.T) {
	config := &config{}
	project := &project{}
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
				err := w.Create("front", project.getPath(t))
				assert.NoError(t, err)
				err = w.CreateEnv("front", "prod")
				assert.NoError(t, err)

				functionPath := config.getPath(t) + "/front/functions/functions.bash"
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
				assert.Len(t, w.Functions.Functions, 2)
				assert.Equal(t,
					Workspace{
						Name: "front",
						Functions: Functions{
							file: fmt.Sprintf("%s/front/functions/functions.bash", config.getPath(t)),
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
						},
						Envs: []Env{
							{
								Name: "default",
								file: fmt.Sprintf("%s/front/envs/default.bash", config.getPath(t)),
							},
							{
								Name: "prod",
								file: fmt.Sprintf("%s/front/envs/prod.bash", config.getPath(t)),
							},
						},
						Config: map[string]string{
							"app":  "bash",
							"path": project.getPath(t),
						},
						dir: fmt.Sprintf("%s/front", config.getPath(t)),
					}, w)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(config.getPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(config.getPath(t)))
			assert.NoError(t, err)
			s.setup(t, w)
			workspace, err := w.Get("front")
			s.test(t, workspace, err)
		})
	}
}

func TestCreate(t *testing.T) {
	config := &config{}
	project := &project{}
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
				path := config.getPath(t)
				envFile, err := os.Stat(path + "/test/envs/default.bash")
				assert.NoError(t, err)
				assert.Equal(t, "default.bash", envFile.Name())
				functionFile, err := os.Stat(path + "/test/functions/functions.bash")
				assert.NoError(t, err)
				assert.Equal(t, "functions.bash", functionFile.Name())
				configFile, err := os.Stat(path + "/test/config.toml")
				assert.NoError(t, err)
				assert.Equal(t, "config.toml", configFile.Name())
				b, err := os.ReadFile(path + "/test/config.toml")
				assert.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("app = 'bash'\npath = '%s'\n", project.getPath(t)), string(b))
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(config.getPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(config.getPath(t)))
			assert.NoError(t, err)
			s.setup(t, w)
			assert.NoError(t, err)
			s.test(t, w.Create("test", project.getPath(t)))
		})
	}
}

func TestCreateEnv(t *testing.T) {
	config := &config{}
	project := &project{}
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
				path := config.getPath(t)
				defaultEnvFile, err := os.Stat(path + "/test/envs/default.bash")
				assert.NoError(t, err)
				assert.Equal(t, "default.bash", defaultEnvFile.Name())
				prodEnvFile, err := os.Stat(path + "/test/envs/prod.bash")
				assert.NoError(t, err)
				assert.Equal(t, "prod.bash", prodEnvFile.Name())
				functionFile, err := os.Stat(path + "/test/functions/functions.bash")
				assert.NoError(t, err)
				assert.Equal(t, "functions.bash", functionFile.Name())
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(config.getPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(config.getPath(t)))
			assert.NoError(t, err)
			s.setup(t, w)
			assert.NoError(t, err)
			assert.NoError(t, w.Create("test", project.getPath(t)))
			s.test(t, w.CreateEnv("test", "prod"))
		})
	}
}

func TestEdit(t *testing.T) {
	config := &config{}
	project := &project{}
	type scenario struct {
		name  string
		setup func(*testing.T, WorkspaceManager, *MockCommander)
		test  func(*testing.T)
	}
	scenarios := []scenario{
		{
			"Edit workspace",
			func(t *testing.T, w WorkspaceManager, exec *MockCommander) {
				exec.On("command", "", "-c", fmt.Sprintf("emacs %s/test/functions/functions.bash", config.getPath(t))).Return(nil)
				err := w.Create("test", project.getPath(t))
				assert.NoError(t, err)
			},
			func(t *testing.T) {
				f, err := os.Stat(fmt.Sprintf("%s/test/envs/default.bash", config.getPath(t)))
				assert.NoError(t, err)
				assert.Equal(t, "default.bash", f.Name())
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(config.getPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(config.getPath(t)))
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
	config := &config{}
	project := &project{}
	type scenario struct {
		name  string
		env   string
		setup func(*testing.T, *MockCommander)
	}
	scenarios := []scenario{
		{
			"Edit default workspace",
			"default",
			func(t *testing.T, exec *MockCommander) {
				exec.On("command", "", "-c", fmt.Sprintf("emacs %s/test/envs/default.bash", config.getPath(t))).Return(nil)
			},
		},
		{
			"Edit prod workspace",
			"prod",
			func(t *testing.T, exec *MockCommander) {
				exec.On("command", "", "-c", fmt.Sprintf("emacs %s/test/envs/prod.bash", config.getPath(t))).Return(nil)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(config.getPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(config.getPath(t)))
			assert.NoError(t, err)
			err = w.Create("test", project.getPath(t))
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
	config := &config{}
	project := &project{}
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
			"default",
			"/bin/bash",
			func(t *testing.T, exec *MockCommander) {
				functionPath := config.getPath(t) + "/test/functions/functions.bash"
				assert.NoError(t, os.WriteFile(functionPath, []byte(`
run-db() {

}
`), 0o777))

				exec.On("command", project.getPath(t), "-c", fmt.Sprintf("export WO_NAME=test && export WO_ENV=default && source %s/test/envs/default.bash && source %s/test/functions/functions.bash && run-db", config.getPath(t), config.getPath(t))).Return(nil)
			},
		},
		{
			"Run a function with a fish shell",
			[]string{"run-db"},
			"default",
			"/bin/fish",
			func(t *testing.T, exec *MockCommander) {
				functionPath := config.getPath(t) + "/test/functions/functions.fish"
				assert.NoError(t, os.WriteFile(functionPath, []byte(`
function run-db

end
`), 0o777))
				exec.On("command", project.getPath(t), "-C", "set -x -g WO_NAME test", "-C", "set -x -g WO_ENV default", "-C", fmt.Sprintf("source %s/test/envs/default.fish", config.getPath(t)), "-C", fmt.Sprintf("source %s/test/functions/functions.fish", config.getPath(t)), "-c", "run-db").Return(nil)
			},
		},
		{
			"Run a function with an env defined and a bash shell",
			[]string{"run-db", "watch"},
			"prod",
			"/bin/bash",
			func(t *testing.T, exec *MockCommander) {
				functionPath := config.getPath(t) + "/test/functions/functions.bash"
				assert.NoError(t, os.WriteFile(functionPath, []byte(`
run-db() {

}
`), 0o777))

				exec.On("command", project.getPath(t), "-c", fmt.Sprintf("export WO_NAME=test && export WO_ENV=prod && source %s/test/envs/prod.bash && source %s/test/functions/functions.bash && run-db watch", config.getPath(t), config.getPath(t))).Return(nil)
			},
		},
		{
			"Run a function with an env defined and a fish shell",
			[]string{"run-db", "watch"},
			"prod",
			"/bin/fish",
			func(t *testing.T, exec *MockCommander) {
				functionPath := config.getPath(t) + "/test/functions/functions.fish"
				assert.NoError(t, os.WriteFile(functionPath, []byte(`
function run-db
end
`), 0o777))
				exec.On("command", project.getPath(t), "-C", "set -x -g WO_NAME test", "-C", "set -x -g WO_ENV prod", "-C", fmt.Sprintf("source %s/test/envs/prod.fish", config.getPath(t)), "-C", fmt.Sprintf("source %s/test/functions/functions.fish", config.getPath(t)), "-c", "run-db watch").Return(nil)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(config.getPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath(s.shell), WithConfigPath(config.getPath(t)))
			assert.NoError(t, err)
			err = w.Create("test", project.getPath(t))
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
	config := &config{}
	project := &project{}
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
				path := config.getPath(t)
				_, err := os.Stat(path + "/test/envs/prod.bash")
				assert.True(t, os.IsNotExist(err))
				_, err = os.Stat(path + "/test/envs/dev.bash")
				assert.True(t, os.IsNotExist(err))
				_, err = os.Stat(path + "/test/functions/functions.bash")
				assert.True(t, os.IsNotExist(err))
				_, err = os.Stat(path + "/test/config.toml")
				assert.True(t, os.IsNotExist(err))

				_, err = os.Stat(path + "/front/functions/functions.bash")
				assert.NoError(t, err)
				_, err = os.Stat(path + "/front/envs/dev.bash")
				assert.NoError(t, err)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(config.getPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(config.getPath(t)))
			assert.NoError(t, err)
			err = w.Create("test", project.getPath(t))
			assert.NoError(t, err)
			err = w.CreateEnv("test", "prod")
			assert.NoError(t, err)
			err = w.CreateEnv("test", "dev")
			assert.NoError(t, err)
			err = w.Create("front", project.getPath(t))
			assert.NoError(t, err)
			err = w.CreateEnv("front", "dev")
			assert.NoError(t, err)
			s.test(t, w.Remove(s.workspace))
		})
	}
}

func TestSetConfig(t *testing.T) {
	config := &config{}
	project := &project{}
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
			"/tmp",
			func(t *testing.T, err error) {
				assert.NoError(t, err)
				b, err := os.ReadFile(fmt.Sprintf("%s/%s", config.getPath(t), "test/config.toml"))
				assert.NoError(t, err)
				assert.Equal(t, []byte("app = 'bash'\npath = '/tmp'\n"), b)
			},
		},
		{
			"Set an unsupported value",
			"test",
			"whetevr",
			"/home/user/project",
			func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			"Set an unexisting path",
			"test",
			"path",
			"/home/user/project",
			func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, `path "/home/user/project" does not exist`)
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
			os.RemoveAll(config.getPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(config.getPath(t)))
			assert.NoError(t, err)
			err = w.Create("test", project.getPath(t))
			assert.NoError(t, err)
			s.test(t, w.SetConfig(s.workspace, map[string]string{s.key: s.value}))
		})
	}
}

func TestBuildAliases(t *testing.T) {
	config := &config{}
	project := &project{}
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
					fmt.Sprintf(`alias c_front="cd %s/front"`, project.getPath(t)),
					fmt.Sprintf(`alias c_test="cd %s/test"`, project.getPath(t)),
				}, aliases)
			},
		},
		{
			"Build aliases with a different prefix",
			"xx",
			func(t *testing.T, aliases []string, e error) {
				assert.NoError(t, e)
				assert.Equal(t, []string{
					fmt.Sprintf(`alias xxfront="cd %s/front"`, project.getPath(t)),
					fmt.Sprintf(`alias xxtest="cd %s/test"`, project.getPath(t)),
				}, aliases)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(config.getPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(config.getPath(t)))
			assert.NoError(t, err)
			testProjectPath := fmt.Sprintf("%s/test", project.getPath(t))
			assert.NoError(t, os.MkdirAll(testProjectPath, 0o777))
			err = w.Create("test", testProjectPath)
			assert.NoError(t, err)
			frontProjectPath := fmt.Sprintf("%s/front", project.getPath(t))
			assert.NoError(t, os.MkdirAll(frontProjectPath, 0o777))
			err = w.Create("front", frontProjectPath)
			assert.NoError(t, err)
			aliases, err := w.BuildAliases(s.prefix)
			s.test(t, aliases, err)
		})
	}
}

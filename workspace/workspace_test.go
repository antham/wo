package workspace

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var configPath string

func getConfigPath(t *testing.T) string {
	if configPath == "" {
		path, err := os.MkdirTemp("/tmp", "wo")
		assert.NoError(t, err)
		configPath = path
	}
	return configPath
}

func TestNewWorkspaceManager(t *testing.T) {
	type scenario struct {
		name  string
		setup func() (string, string, string)
		test  func(WorkspaceManager, error)
	}
	scenarios := []scenario{
		{
			"No variable editor, nor visual defined",
			func() (string, string, string) {
				return "", "", ""
			},
			func(w WorkspaceManager, err error) {
				assert.EqualError(t, err, "no editor defined")
			},
		},
		{
			"visual defined",
			func() (string, string, string) {
				return "vim", "", "/usr/bin/bash"
			},
			func(w WorkspaceManager, err error) {
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
			func(w WorkspaceManager, err error) {
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
			s.test(NewWorkspaceManager(WithEditor(editor, visual), WithShellPath(shell), WithConfigPath(getConfigPath(t))))
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
				assert.NoError(t, w.Create("api"))
				assert.NoError(t, w.CreateEnv("api", "dev"))
				assert.NoError(t, w.Create("db"))
				assert.NoError(t, w.CreateEnv("db", "staging"))
				assert.NoError(t, w.Create("front"))
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
		test  func(Workspace, error)
	}
	scenarios := []scenario{
		{
			"Workspace does not exist",
			func(t *testing.T, w WorkspaceManager) {
			},
			func(w Workspace, err error) {
				assert.EqualError(t, err, "the workspace does not exist")
			},
		},
		{
			"Get all workspace",
			func(t *testing.T, w WorkspaceManager) {
				err := w.Create("front")
				assert.NoError(t, err)
				err = w.CreateEnv("front", "prod")
				assert.NoError(t, err)

				functionPath := getConfigPath(t) + "/functions/bash/front.bash"
				assert.NoError(t, os.WriteFile(functionPath, []byte(`
# A function 1
test_func1() {

}

# A function 2
test_func2() {

}
`), 0o777))
			},
			func(w Workspace, err error) {
				assert.NoError(t, err)
				assert.Len(t, w.Functions, 2)
				assert.Equal(t, Workspace{
					Name: "front",
					Functions: []Function{
						{
							Function:    "test_func1",
							Description: "A function 1",
						},
						{
							Function:    "test_func2",
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
			s.test(w.Get("front"))
		})
	}
}

func TestCreate(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T, WorkspaceManager)
		test  func(error)
	}
	scenarios := []scenario{
		{
			"Create workspace",
			func(t *testing.T, w WorkspaceManager) {
			},
			func(err error) {
				assert.NoError(t, err)
				path := getConfigPath(t)
				envFile, err := os.Stat(path + "/envs/bash/test/default.bash")
				assert.Equal(t, "default.bash", envFile.Name())
				assert.NoError(t, err)
				functionFile, err := os.Stat(path + "/functions/bash/test.bash")
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
			s.test(w.Create("test"))
		})
	}
}

func TestCreateEnv(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T, WorkspaceManager)
		test  func(error)
	}
	scenarios := []scenario{
		{
			"Create workspace env",
			func(t *testing.T, w WorkspaceManager) {
			},
			func(err error) {
				assert.NoError(t, err)
				path := getConfigPath(t)
				defaultEnvFile, err := os.Stat(path + "/envs/bash/test/default.bash")
				assert.Equal(t, "default.bash", defaultEnvFile.Name())
				assert.NoError(t, err)
				prodEnvFile, err := os.Stat(path + "/envs/bash/test/prod.bash")
				assert.Equal(t, "prod.bash", prodEnvFile.Name())
				assert.NoError(t, err)
				functionFile, err := os.Stat(path + "/functions/bash/test.bash")
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
			assert.NoError(t, w.Create("test"))
			s.test(w.CreateEnv("test", "prod"))
		})
	}
}

func TestEdit(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T, WorkspaceManager)
		test  func([]string)
	}
	scenarios := []scenario{
		{
			"Edit workspace",
			func(t *testing.T, w WorkspaceManager) {
				err := w.Create("test")
				assert.NoError(t, err)
			},
			func(args []string) {
				assert.Equal(t, []string{"-c", fmt.Sprintf("emacs %s/functions/bash/test.bash", getConfigPath(t))}, args)
				f, err := os.Stat(fmt.Sprintf("%s/envs/bash/test/default.bash", getConfigPath(t)))
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
			s.setup(t, w)
			args := []string{}
			w.execCommand = func(a ...string) error {
				args = a
				return nil
			}
			assert.NoError(t, err)
			assert.NoError(t, w.Edit("test"))
			s.test(args)
		})
	}
}

func TestEditEnv(t *testing.T) {
	type scenario struct {
		name string
		env  string
		test func([]string)
	}
	scenarios := []scenario{
		{
			"Edit default workspace",
			"",
			func(args []string) {
				assert.Equal(t, []string{"-c", fmt.Sprintf("emacs %s/envs/bash/test/default.bash", getConfigPath(t))}, args)
			},
		},
		{
			"Edit prod workspace",
			"prod",
			func(args []string) {
				assert.Equal(t, []string{"-c", fmt.Sprintf("emacs %s/envs/bash/test/prod.bash", getConfigPath(t))}, args)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(getConfigPath(t)))
			assert.NoError(t, err)
			err = w.Create("test")
			assert.NoError(t, err)
			err = w.CreateEnv("test", "prod")
			assert.NoError(t, err)
			args := []string{}
			w.execCommand = func(a ...string) error {
				args = a
				return nil
			}
			assert.NoError(t, err)
			assert.NoError(t, w.EditEnv("test", s.env))
			s.test(args)
		})
	}
}

func TestLoad(t *testing.T) {
	type scenario struct {
		name  string
		env   string
		shell string
		test  func([]string)
	}
	scenarios := []scenario{
		{
			"Load workspace with a bash shell",
			"",
			"/bin/bash",
			func(args []string) {
				assert.Equal(t, []string{"-c", fmt.Sprintf("export WO_NAME=test && export WO_ENV=default && source %s/envs/bash/test/default.bash && source %s/functions/bash/test.bash", getConfigPath(t), getConfigPath(t))}, args)
			},
		},
		{
			"Load workspace with a fish shell",
			"",
			"/bin/fish",
			func(args []string) {
				assert.Equal(t, []string{"-C", "set -x -g WO_NAME test", "-C", "set -x -g WO_ENV default", "-C", fmt.Sprintf("source %s/envs/fish/test/default.fish", getConfigPath(t)), "-C", fmt.Sprintf("source %s/functions/fish/test.fish", getConfigPath(t))}, args)
			},
		},
		{
			"Load workspace with an env defined and a bash shell",
			"prod",
			"/bin/bash",
			func(args []string) {
				assert.Equal(t, []string{"-c", fmt.Sprintf("export WO_NAME=test && export WO_ENV=prod && source %s/envs/bash/test/prod.bash && source %s/functions/bash/test.bash", getConfigPath(t), getConfigPath(t))}, args)
			},
		},
		{
			"Load workspace with an env defined and a fish shell",
			"prod",
			"/bin/fish",
			func(args []string) {
				assert.Equal(t, []string{"-C", "set -x -g WO_NAME test", "-C", "set -x -g WO_ENV prod", "-C", fmt.Sprintf("source %s/envs/fish/test/prod.fish", getConfigPath(t)), "-C", fmt.Sprintf("source %s/functions/fish/test.fish", getConfigPath(t))}, args)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath(s.shell), WithConfigPath(getConfigPath(t)))
			assert.NoError(t, err)
			err = w.Create("test")
			assert.NoError(t, err)
			err = w.CreateEnv("test", "prod")
			assert.NoError(t, err)
			args := []string{}
			w.execCommand = func(a ...string) error {
				args = a
				return nil
			}
			assert.NoError(t, err)
			assert.NoError(t, w.Load("test", s.env))
			s.test(args)
		})
	}
}

func TestRunFunction(t *testing.T) {
	type scenario struct {
		name            string
		functionAndArgs []string
		env             string
		shell           string
		test            func([]string)
	}
	scenarios := []scenario{
		{
			"Run a function with a bash shell",
			[]string{"run-db"},
			"",
			"/bin/bash",
			func(args []string) {
				assert.Equal(t, []string{"-c", fmt.Sprintf("export WO_NAME=test && export WO_ENV=default && source %s/envs/bash/test/default.bash && source %s/functions/bash/test.bash && run-db", getConfigPath(t), getConfigPath(t))}, args)
			},
		},
		{
			"Run a function with a fish shell",
			[]string{"run-db"},
			"",
			"/bin/fish",
			func(args []string) {
				assert.Equal(t, []string{"-C", "set -x -g WO_NAME test", "-C", "set -x -g WO_ENV default", "-C", fmt.Sprintf("source %s/envs/fish/test/default.fish", getConfigPath(t)), "-C", fmt.Sprintf("source %s/functions/fish/test.fish", getConfigPath(t)), "-c", "run-db"}, args)
			},
		},
		{
			"Run a function with an env defined and a bash shell",
			[]string{"run-db", "watch"},
			"prod",
			"/bin/bash",
			func(args []string) {
				assert.Equal(t, []string{"-c", fmt.Sprintf("export WO_NAME=test && export WO_ENV=prod && source "+getConfigPath(t)+"/envs/bash/test/prod.bash && source %s/functions/bash/test.bash && run-db watch", getConfigPath(t))}, args)
			},
		},
		{
			"Run a function with an env defined and a fish shell",
			[]string{"run-db", "watch"},
			"prod",
			"/bin/fish",
			func(args []string) {
				assert.Equal(t, []string{"-C", "set -x -g WO_NAME test", "-C", "set -x -g WO_ENV prod", "-C", fmt.Sprintf("source %s/envs/fish/test/prod.fish", getConfigPath(t)), "-C", fmt.Sprintf("source %s/functions/fish/test.fish", getConfigPath(t)), "-c", "run-db watch"}, args)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath(s.shell), WithConfigPath(getConfigPath(t)))
			assert.NoError(t, err)
			err = w.Create("test")
			assert.NoError(t, err)
			err = w.CreateEnv("test", "prod")
			assert.NoError(t, err)
			args := []string{}
			w.execCommand = func(a ...string) error {
				args = a
				return nil
			}
			assert.NoError(t, err)
			assert.NoError(t, w.RunFunction("test", s.env, s.functionAndArgs))
			s.test(args)
		})
	}
}

func TestRemove(t *testing.T) {
	type scenario struct {
		name      string
		workspace string
		test      func(error)
	}
	scenarios := []scenario{
		{
			"Remove an unexisting workspace",
			"whatever",
			func(e error) {
				assert.Error(t, e)
			},
		},
		{
			"Remove a workspace",
			"test",
			func(e error) {
				assert.NoError(t, e)
				path := getConfigPath(t)
				_, err := os.Stat(path + "/envs/bash/test/prod.bash")
				assert.True(t, os.IsNotExist(err))
				_, err = os.Stat(path + "/envs/bash/test/dev.bash")
				assert.True(t, os.IsNotExist(err))
				_, err = os.Stat(path + "/functions/bash/test.bash")
				assert.True(t, os.IsNotExist(err))

				_, err = os.Stat(path + "/functions/bash/front.bash")
				assert.NoError(t, err)
				_, err = os.Stat(path + "/envs/bash/front/dev.bash")
				assert.NoError(t, err)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			w, err := NewWorkspaceManager(WithEditor("emacs", "emacs"), WithShellPath("/bin/bash"), WithConfigPath(getConfigPath(t)))
			assert.NoError(t, err)
			err = w.Create("test")
			assert.NoError(t, err)
			err = w.CreateEnv("test", "prod")
			assert.NoError(t, err)
			err = w.CreateEnv("test", "dev")
			assert.NoError(t, err)
			err = w.Create("front")
			assert.NoError(t, err)
			err = w.CreateEnv("front", "dev")
			assert.NoError(t, err)
			s.test(w.Remove(s.workspace))
		})
	}
}

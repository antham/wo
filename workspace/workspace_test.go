package workspace

import (
	"os"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getConfigPath(t *testing.T) string {
	usr, err := user.Current()
	assert.NoError(t, err)
	homeDir := usr.HomeDir
	return homeDir + "/" + configDir
}

func TestNewWorkspaceManager(t *testing.T) {
	type scenario struct {
		name  string
		setup func()
		test  func(WorkspaceManager, error)
	}
	scenarios := []scenario{
		{
			"No variable EDITOR, nor VISUAL defined",
			func() {
			},
			func(w WorkspaceManager, err error) {
				assert.EqualError(t, err, "no VISUAL or EDITOR environment variable found")
			},
		},
		{
			"VISUAL defined",
			func() {
				os.Setenv("VISUAL", "vim")
				os.Setenv("SHELL", "bash")
			},
			func(w WorkspaceManager, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "vim", w.editor)
				assert.Equal(t, "bash", w.shellBin)
				assert.DirExists(t, w.getFunctionDir())
				assert.DirExists(t, w.getConfigDir())
				assert.DirExists(t, w.getEnvDir())
			},
		},
		{
			"EDITOR defined",
			func() {
				os.Setenv("EDITOR", "emacs")
				os.Setenv("SHELL", "zsh")
			},
			func(w WorkspaceManager, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "emacs", w.editor)
				assert.Equal(t, "zsh", w.shellBin)
				assert.DirExists(t, w.getFunctionDir())
				assert.DirExists(t, w.getConfigDir())
				assert.DirExists(t, w.getEnvDir())
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.Clearenv()
			os.RemoveAll(getConfigPath(t))
			s.setup()
			s.test(NewWorkspaceManager())
		})
	}

}

func TestList(t *testing.T) {
	type scenario struct {
		name  string
		setup func()
		test  func([]Workspace, error)
	}
	scenarios := []scenario{
		{
			"Invalid function file",
			func() {
				path := getConfigPath(t) + "/functions/bash"
				assert.NoError(t, os.WriteFile(path+"/front", []byte{}, 0777))
			},
			func(ws []Workspace, err error) {
				assert.ErrorContains(t, err, "the workspace does not exist")
			},
		},
		{
			"No workspace defined",
			func() {
			},
			func(ws []Workspace, err error) {
				assert.NoError(t, err)
				assert.Empty(t, ws)
			},
		},
		{
			"Get all workspaces ordered alphabetically",
			func() {
				functionPath := getConfigPath(t) + "/functions/bash"
				assert.NoError(t, os.WriteFile(functionPath+"/front.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(functionPath+"/api.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(functionPath+"/db.bash", []byte{}, 0777))

				envPath := getConfigPath(t) + "/envs/bash"
				assert.NoError(t, os.MkdirAll(envPath+"/front", 0777))
				assert.NoError(t, os.MkdirAll(envPath+"/api", 0777))
				assert.NoError(t, os.MkdirAll(envPath+"/db", 0777))
				assert.NoError(t, os.WriteFile(envPath+"/front/prod.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(envPath+"/api/dev.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(envPath+"/db/staging.bash", []byte{}, 0777))
			},
			func(ws []Workspace, err error) {
				assert.NoError(t, err)
				assert.Len(t, ws, 3)
				assert.Equal(t, []Workspace{
					{Name: "api", Functions: []Function{}, Envs: []string{"dev"}},
					{Name: "db", Functions: []Function{}, Envs: []string{"staging"}},
					{Name: "front", Functions: []Function{}, Envs: []string{"prod"}},
				}, ws)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			os.Setenv("VISUAL", "emacs")
			os.Setenv("EDITOR", "emacs")
			os.Setenv("SHELL", "bash")
			w, err := NewWorkspaceManager()
			assert.NoError(t, err)
			s.setup()
			s.test(w.List())
		})
	}
}

func TestGet(t *testing.T) {
	type scenario struct {
		name  string
		setup func()
		test  func(Workspace, error)
	}
	scenarios := []scenario{
		{
			"Workspace does not exist",
			func() {
			},
			func(w Workspace, err error) {
				assert.EqualError(t, err, "the workspace does not exist")
			},
		},
		{
			"Get all workspace",
			func() {
				functionPath := getConfigPath(t) + "/functions/bash"
				assert.NoError(t, os.WriteFile(functionPath+"/front.bash", []byte(`
# A function 1
test_func1() {

}

# A function 2
test_func2() {

}
`), 0777))

				envPath := getConfigPath(t) + "/envs/bash"
				assert.NoError(t, os.MkdirAll(envPath+"/front", 0777))
				assert.NoError(t, os.WriteFile(envPath+"/front/prod.bash", []byte(``), 0777))
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
					Envs: []string{"prod"},
				}, w)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			os.Setenv("VISUAL", "emacs")
			os.Setenv("EDITOR", "emacs")
			os.Setenv("SHELL", "bash")
			w, err := NewWorkspaceManager()
			assert.NoError(t, err)
			s.setup()
			s.test(w.Get("front"))
		})
	}

}

func TestEdit(t *testing.T) {
	type scenario struct {
		name string
		test func([]string)
	}
	scenarios := []scenario{
		{
			"Edit workspace",
			func(args []string) {
				assert.Equal(t, []string{"-c", "emacs " + getHomePath(t) + "/.config/wo/functions/bash/test.bash"}, args)
				f, err := os.Stat(getHomePath(t) + "/.config/wo/envs/bash/test/default.bash")
				assert.NoError(t, err)
				assert.Equal(t, "default.bash", f.Name())
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			os.Setenv("VISUAL", "emacs")
			os.Setenv("EDITOR", "emacs")
			os.Setenv("SHELL", "bash")
			w, err := NewWorkspaceManager()
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
				assert.Equal(t, []string{"-c", "emacs " + getHomePath(t) + "/.config/wo/envs/bash/test/default.bash"}, args)
			},
		},
		{
			"Edit prod workspace",
			"prod",
			func(args []string) {
				assert.Equal(t, []string{"-c", "emacs " + getHomePath(t) + "/.config/wo/envs/bash/test/prod.bash"}, args)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			os.Setenv("VISUAL", "emacs")
			os.Setenv("EDITOR", "emacs")
			os.Setenv("SHELL", "bash")
			w, err := NewWorkspaceManager()
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
		setup func()
		test  func([]string)
	}
	scenarios := []scenario{
		{
			"Load workspace with a bash shell",
			"",
			func() {
				os.Setenv("SHELL", "bash")
			},
			func(args []string) {
				assert.Equal(t, []string{"-c", "export WO_NAME=test && export WO_ENV=default && source " + getHomePath(t) + "/.config/wo/functions/bash/test.bash"}, args)
			},
		},
		{
			"Load workspace with a fish shell",
			"",
			func() {
				os.Setenv("SHELL", "fish")
			},
			func(args []string) {
				assert.Equal(t, []string{"-C", "set -x -g WO_NAME test", "-C", "set -x -g WO_ENV default", "-C", "source " + getHomePath(t) + "/.config/wo/functions/fish/test.fish"}, args)
			},
		},
		{
			"Load workspace with an env defined and a bash shell",
			"prod",
			func() {
				envPath := getConfigPath(t) + "/envs/bash"
				assert.NoError(t, os.MkdirAll(envPath+"/test", 0777))
				assert.NoError(t, os.WriteFile(envPath+"/test/prod.bash", []byte{}, 0777))
				os.Setenv("SHELL", "bash")
			},
			func(args []string) {
				assert.Equal(t, []string{"-c", "export WO_NAME=test && export WO_ENV=prod && source " + getHomePath(t) + "/.config/wo/envs/bash/test/prod.bash && source " + getHomePath(t) + "/.config/wo/functions/bash/test.bash"}, args)
			},
		},
		{
			"Load workspace with an env defined and a fish shell",
			"prod",
			func() {
				envPath := getConfigPath(t) + "/envs/fish"
				assert.NoError(t, os.MkdirAll(envPath+"/test", 0777))
				assert.NoError(t, os.WriteFile(envPath+"/test/prod.fish", []byte{}, 0777))
				os.Setenv("SHELL", "fish")
			},
			func(args []string) {
				assert.Equal(t, []string{"-C", "set -x -g WO_NAME test", "-C", "set -x -g WO_ENV prod", "-C", "source " + getHomePath(t) + "/.config/wo/envs/fish/test/prod.fish", "-C", "source " + getHomePath(t) + "/.config/wo/functions/fish/test.fish"}, args)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			os.Setenv("VISUAL", "emacs")
			os.Setenv("EDITOR", "emacs")
			os.Unsetenv("SHELL")
			s.setup()
			w, err := NewWorkspaceManager()
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
		setup           func()
		test            func([]string)
	}
	scenarios := []scenario{
		{
			"Run a function with a bash shell",
			[]string{"run-db"},
			"",
			func() {
				os.Setenv("SHELL", "/bin/bash")
			},
			func(args []string) {
				assert.Equal(t, []string{"-c", "export WO_NAME=test && export WO_ENV=default && source " + getHomePath(t) + "/.config/wo/functions/bash/test.bash && run-db"}, args)
			},
		},
		{
			"Run a function with a fish shell",
			[]string{"run-db"},
			"",
			func() {
				os.Setenv("SHELL", "/bin/fish")
			},
			func(args []string) {
				assert.Equal(t, []string{"-C", "set -x -g WO_NAME test", "-C", "set -x -g WO_ENV default", "-C", "source " + getHomePath(t) + "/.config/wo/functions/fish/test.fish", "-c", "run-db"}, args)
			},
		},
		{
			"Run a function with an env defined and a bash shell",
			[]string{"run-db", "watch"},
			"prod",
			func() {
				envPath := getConfigPath(t) + "/envs/bash"
				assert.NoError(t, os.MkdirAll(envPath+"/test", 0777))
				assert.NoError(t, os.WriteFile(envPath+"/test/prod.bash", []byte{}, 0777))
				os.Setenv("SHELL", "/bin/bash")
			},
			func(args []string) {
				assert.Equal(t, []string{"-c", "export WO_NAME=test && export WO_ENV=prod && source " + getHomePath(t) + "/.config/wo/envs/bash/test/prod.bash && source " + getHomePath(t) + "/.config/wo/functions/bash/test.bash && run-db watch"}, args)
			},
		},
		{
			"Run a function with an env defined and a fish shell",
			[]string{"run-db", "watch"},
			"prod",
			func() {
				envPath := getConfigPath(t) + "/envs/fish"
				assert.NoError(t, os.MkdirAll(envPath+"/test", 0777))
				assert.NoError(t, os.WriteFile(envPath+"/test/prod.fish", []byte{}, 0777))
				os.Setenv("SHELL", "/bin/fish")
			},
			func(args []string) {
				assert.Equal(t, []string{"-C", "set -x -g WO_NAME test", "-C", "set -x -g WO_ENV prod", "-C", "source " + getHomePath(t) + "/.config/wo/envs/fish/test/prod.fish", "-C", "source " + getHomePath(t) + "/.config/wo/functions/fish/test.fish", "-c", "run-db watch"}, args)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			os.Unsetenv("SHELL")
			os.Setenv("VISUAL", "emacs")
			os.Setenv("EDITOR", "emacs")
			s.setup()
			w, err := NewWorkspaceManager()
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
		setup     func()
		test      func(error)
	}
	scenarios := []scenario{
		{
			"Remove an unexisting workspace",
			"whatever",
			func() {
			},
			func(e error) {
				assert.Error(t, e)
			},
		},
		{
			"Remove a workspace",
			"test",
			func() {
				path := getConfigPath(t)
				assert.NoError(t, os.MkdirAll(path+"/envs/bash/test", 0777))
				assert.NoError(t, os.MkdirAll(path+"/envs/bash/front", 0777))
				assert.NoError(t, os.WriteFile(path+"/envs/bash/test/prod.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(path+"/envs/bash/test/dev.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(path+"/envs/bash/front/dev.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(path+"/functions/bash/front.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(path+"/functions/bash/test.bash", []byte{}, 0777))
			},
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
			os.Setenv("VISUAL", "emacs")
			os.Setenv("EDITOR", "emacs")
			os.Setenv("SHELL", "bash")
			w, err := NewWorkspaceManager()
			assert.NoError(t, err)
			s.setup()
			s.test(w.Remove("test"))
		})
	}
}

func getHomePath(t *testing.T) string {
	u, err := user.Current()
	assert.NoError(t, err)
	return u.HomeDir
}

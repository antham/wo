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
				assert.Len(t, ws, 0)
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
				assert.NoError(t, os.WriteFile(envPath+"/front.prod.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(envPath+"/api.dev.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(envPath+"/db.staging.bash", []byte{}, 0777))
			},
			func(ws []Workspace, err error) {
				assert.NoError(t, err)
				assert.Len(t, ws, 3)
				assert.Equal(t, []Workspace{
					{Name: "api", Commands: []Command{}, Envs: []string{"dev"}},
					{Name: "db", Commands: []Command{}, Envs: []string{"staging"}},
					{Name: "front", Commands: []Command{}, Envs: []string{"prod"}},
				}, ws)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			os.Setenv("VISUAL", "emacs")
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
				assert.NoError(t, os.WriteFile(envPath+"/front.prod.bash", []byte(``), 0777))
			},
			func(w Workspace, err error) {
				assert.NoError(t, err)
				assert.Len(t, w.Commands, 2)
				assert.Equal(t, Workspace{
					Name: "front",
					Commands: []Command{
						{
							Command:     "test_func1",
							Description: "A function 1",
						},
						{
							Command:     "test_func2",
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
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			os.Setenv("VISUAL", "emacs")
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
				assert.Equal(t, []string{"-c", "emacs " + getHomePath(t) + "/.config/wo/envs/bash/test.default.bash"}, args)
			},
		},
		{
			"Edit prod workspace",
			"prod",
			func(args []string) {
				assert.Equal(t, []string{"-c", "emacs " + getHomePath(t) + "/.config/wo/envs/bash/test.prod.bash"}, args)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			os.Setenv("VISUAL", "emacs")
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
			"Load workspace",
			"",
			func() {
			},
			func(args []string) {
				assert.Equal(t, []string{"-C", "source " + getHomePath(t) + "/.config/wo/functions/bash/test.bash"}, args)
			},
		},
		{
			"Load workspace with an env defined",
			"prod",
			func() {
				path := getConfigPath(t) + "/envs/bash"
				assert.NoError(t, os.WriteFile(path+"/test.prod.bash", []byte{}, 0777))
			},
			func(args []string) {
				assert.Equal(t, []string{"-C", "source " + getHomePath(t) + "/.config/wo/envs/bash/test.prod.bash", "-C", "source " + getHomePath(t) + "/.config/wo/functions/bash/test.bash"}, args)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			os.Setenv("VISUAL", "emacs")
			os.Setenv("SHELL", "bash")
			w, err := NewWorkspaceManager()
			args := []string{}
			w.execCommand = func(a ...string) error {
				args = a
				return nil
			}
			assert.NoError(t, err)
			s.setup()
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
			"Run function",
			[]string{"run-db"},
			"",
			func() {
			},
			func(args []string) {
				assert.Equal(t, []string{"-C", "source " + getHomePath(t) + "/.config/wo/functions/bash/test.bash", "-c", "run-db"}, args)
			},
		},
		{
			"Run a function with an env defined",
			[]string{"run-db", "--watch"},
			"prod",
			func() {
				path := getConfigPath(t) + "/envs/bash"
				assert.NoError(t, os.WriteFile(path+"/test.prod.bash", []byte{}, 0777))
			},
			func(args []string) {
				assert.Equal(t, []string{"-C", "source " + getHomePath(t) + "/.config/wo/envs/bash/test.prod.bash", "-C", "source " + getHomePath(t) + "/.config/wo/functions/bash/test.bash", "-c", "run-db --watch"}, args)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			os.Setenv("VISUAL", "emacs")
			os.Setenv("SHELL", "bash")
			w, err := NewWorkspaceManager()
			args := []string{}
			w.execCommand = func(a ...string) error {
				args = a
				return nil
			}
			assert.NoError(t, err)
			s.setup()
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
				assert.EqualError(t, e, "the workspace does not exist")
			},
		},
		{
			"Remove a workspace",
			"test",
			func() {
				path := getConfigPath(t)
				assert.NoError(t, os.WriteFile(path+"/envs/bash/test.prod.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(path+"/envs/bash/test.dev.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(path+"/envs/bash/front.dev.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(path+"/functions/bash/front.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(path+"/functions/bash/test.bash", []byte{}, 0777))
			},
			func(e error) {
				assert.NoError(t, e)
				path := getConfigPath(t)
				_, err := os.Stat(path + "/envs/bash/test.prod.bash")
				assert.True(t, os.IsNotExist(err))
				_, err = os.Stat(path + "/envs/bash/test.dev.bash")
				assert.True(t, os.IsNotExist(err))
				_, err = os.Stat(path + "/functions/bash/test.bash")
				assert.True(t, os.IsNotExist(err))

				_, err = os.Stat(path + "/functions/bash/front.bash")
				assert.NoError(t, err)
				_, err = os.Stat(path + "/envs/bash/front.dev.bash")
				assert.NoError(t, err)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.RemoveAll(getConfigPath(t))
			os.Setenv("VISUAL", "emacs")
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

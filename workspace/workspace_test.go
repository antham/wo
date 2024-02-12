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
				assert.Equal(t, "bash", w.shell)
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
				assert.Equal(t, "zsh", w.shell)
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
				path := getConfigPath(t) + "/functions/"
				assert.NoError(t, os.WriteFile(path+"/front", []byte{}, 0777))
			},
			func(ws []Workspace, err error) {
				assert.ErrorContains(t, err, "no such file or directory")
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
				path := getConfigPath(t) + "/functions/"
				assert.NoError(t, os.WriteFile(path+"/front.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(path+"/api.bash", []byte{}, 0777))
				assert.NoError(t, os.WriteFile(path+"/db.bash", []byte{}, 0777))
			},
			func(ws []Workspace, err error) {
				assert.NoError(t, err)
				assert.Len(t, ws, 3)
				assert.Equal(t, []Workspace{
					{Name: "api", Commands: []Command{}},
					{Name: "db", Commands: []Command{}},
					{Name: "front", Commands: []Command{}},
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

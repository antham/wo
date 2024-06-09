package cmd

import (
	"fmt"
	"os"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkspaceManager(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T)
		test  func(*testing.T, workspaceManager, error)
	}
	scenarios := []scenario{
		{
			"Missing VISUAL or EDITOR variable",
			func(t *testing.T) {
				os.Unsetenv("VISUAL")
				os.Unsetenv("EDITOR")
			},
			func(t *testing.T, w workspaceManager, err error) {
				assert.Error(t, err)
			},
		},
		{
			"Missing SHELL variable",
			func(t *testing.T) {
				os.Setenv("VISUAL", "emacs")
				os.Unsetenv("SHELL")
			},
			func(t *testing.T, w workspaceManager, err error) {
				assert.Error(t, err)
			},
		},
		{
			"Customize the config directory",
			func(t *testing.T) {
				os.Setenv("VISUAL", "emacs")
				os.Setenv("SHELL", "/bin/bash")
				os.Setenv("WO_CONFIG_PATH", os.TempDir())
			},
			func(t *testing.T, w workspaceManager, err error) {
				assert.NoError(t, err)
				path := os.Getenv("WO_CONFIG_PATH")
				fileInfo, err := os.Stat(path)
				assert.NoError(t, err)
				assert.True(t, fileInfo.IsDir())
			},
		},
		{
			"Create the config in the default folder",
			func(t *testing.T) {
				os.Setenv("VISUAL", "emacs")
				os.Setenv("SHELL", "/bin/bash")
			},
			func(t *testing.T, w workspaceManager, err error) {
				assert.NoError(t, err)
				usr, err := user.Current()
				assert.NoError(t, err)
				fileInfo, err := os.Stat(fmt.Sprintf("%s/.config/wo", usr.HomeDir))
				assert.NoError(t, err)
				assert.True(t, fileInfo.IsDir())
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			s.setup(t)
			w, err := newWorkspaceManager()
			s.test(t, w, err)
		})
	}
}

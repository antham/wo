package cmd

import (
	"os"
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

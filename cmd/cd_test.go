package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCdCmd(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T) (workspaceManager, []string)
		test  func(*testing.T, error)
	}
	scenarios := []scenario{
		{
			"An error occurred when jumping to the folder workspace",
			func(t *testing.T) (workspaceManager, []string) {
				w := newMockWorkspaceManager(t)
				args := []string{"api"}
				w.Mock.On("Cd", args[0]).Return(errors.New("an error occurred"))
				return w, args
			},
			func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			"Jumping to the workspace directory successfully",
			func(t *testing.T) (workspaceManager, []string) {
				w := newMockWorkspaceManager(t)
				args := []string{"api"}
				w.Mock.On("Cd", args[0]).Return(nil)
				return w, args
			},
			func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.Setenv("EDITOR", "emacs")
			os.Setenv("SHELL", "/bin/sh")
			w, args := s.setup(t)
			cmd := newCdCmd(w, newMockCompletionManager(t))
			cmd.SetArgs(args)
			cmd.SetErr(&bytes.Buffer{})
			cmd.SetOut(&bytes.Buffer{})
			s.test(t, cmd.Execute())
		})
	}
}

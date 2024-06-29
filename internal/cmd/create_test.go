package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCreateCmd(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T) (workspaceManager, []string)
		test  func(*testing.T, *bytes.Buffer, *bytes.Buffer, error)
	}
	scenarios := []scenario{
		{
			"An error occurred when creating a workspace",
			func(t *testing.T) (workspaceManager, []string) {
				w := newMockWorkspaceManager(t)
				args := []string{"api", "/tmp/project"}
				w.Mock.On("Create", args[0], args[1]).Return(errors.New("an error occurred"))
				return w, args
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.Error(t, err)
			},
		},
		{
			"Creating an invalid workspace",
			func(t *testing.T) (workspaceManager, []string) {
				return newMockWorkspaceManager(t), []string{"api%", "/tmp/project"}
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.Error(t, err)
			},
		},
		{
			"Creating a workspace successfully",
			func(t *testing.T) (workspaceManager, []string) {
				w := newMockWorkspaceManager(t)
				args := []string{"api", "/tmp/project"}
				w.Mock.On("Create", args[0], args[1]).Return(nil)
				return w, args
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "Workspace 'api' created on path '/tmp/project', reload your shell to setup the project aliases\n", outBuf.String())
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.Setenv("EDITOR", "emacs")
			os.Setenv("SHELL", "/bin/sh")
			errBuf := &bytes.Buffer{}
			outBuf := &bytes.Buffer{}
			w, args := s.setup(t)
			cmd := newCreateCmd(w, newMockCompletionManager(t))
			cmd.SetArgs(args)
			cmd.SetErr(errBuf)
			cmd.SetOut(outBuf)
			s.test(t, outBuf, errBuf, cmd.Execute())
		})
	}
}

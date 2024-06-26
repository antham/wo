package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEditEnvCmd(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T) (workspaceManager, []string)
		test  func(*testing.T, *bytes.Buffer, *bytes.Buffer, error)
	}
	scenarios := []scenario{
		{
			"An error occurred when editing a workspace env",
			func(t *testing.T) (workspaceManager, []string) {
				w := newMockWorkspaceManager(t)
				args := []string{"api", "prod"}
				w.Mock.On("EditEnv", args[0], args[1]).Return(errors.New("an error occurred"))
				return w, args
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.Error(t, err)
			},
		},
		{
			"Editing a workspace env successfully",
			func(t *testing.T) (workspaceManager, []string) {
				w := newMockWorkspaceManager(t)
				args := []string{"api", "prod"}
				w.Mock.On("EditEnv", args[0], args[1]).Return(nil)
				return w, args
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "Environment 'prod' edited on workspace 'api'\n", outBuf.String())
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
			cmd := newEditEnvCmd(w, newMockCompletionManager(t))
			cmd.SetArgs(args)
			cmd.SetErr(errBuf)
			cmd.SetOut(outBuf)
			s.test(t, outBuf, errBuf, cmd.Execute())
		})
	}
}

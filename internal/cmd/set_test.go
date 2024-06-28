package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSetCmd(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T) (workspaceManager, []string)
		test  func(*testing.T, *bytes.Buffer, *bytes.Buffer, error)
	}
	scenarios := []scenario{
		{
			"An error occurred when setting a config",
			func(t *testing.T) (workspaceManager, []string) {
				w := newMockWorkspaceManager(t)
				args := []string{"api", "path", "/home/user/project"}
				w.Mock.On("SetConfig", args[0], map[string]string{args[1]: args[2]}).Return(errors.New("an error occurred"))
				return w, args
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.Error(t, err)
			},
		},
		{
			"Setting a config successfully",
			func(t *testing.T) (workspaceManager, []string) {
				w := newMockWorkspaceManager(t)
				args := []string{"api", "path", "/home/user/project"}
				w.Mock.On("SetConfig", args[0], map[string]string{args[1]: args[2]}).Return(nil)
				return w, args
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "Config key 'path' edited on workspace 'api'\n", outBuf.String())
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
			cmd := newSetCmd(w, newMockCompletionManager(t))
			cmd.SetArgs(args)
			cmd.SetErr(errBuf)
			cmd.SetOut(outBuf)
			s.test(t, outBuf, errBuf, cmd.Execute())
		})
	}
}

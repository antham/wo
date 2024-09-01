package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFixCmd(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T) workspaceManager
		test  func(*testing.T, *bytes.Buffer, *bytes.Buffer, error)
	}
	scenarios := []scenario{
		{
			"An error occurred when fixing the config",
			func(t *testing.T) workspaceManager {
				w := newMockWorkspaceManager(t)
				w.Mock.On("Fix").Return(errors.New("an error occurred"))
				return w
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.Error(t, err)
			},
		},
		{
			"Creating a workspace env successfully",
			func(t *testing.T) workspaceManager {
				w := newMockWorkspaceManager(t)
				w.Mock.On("Fix").Return(nil)
				return w
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "Config folder fixed", outBuf.String())
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			os.Setenv("EDITOR", "emacs")
			os.Setenv("SHELL", "/bin/sh")
			errBuf := &bytes.Buffer{}
			outBuf := &bytes.Buffer{}
			w := s.setup(t)
			cmd := newFixCmd(w)
			cmd.SetArgs([]string{})
			cmd.SetErr(errBuf)
			cmd.SetOut(outBuf)
			s.test(t, outBuf, errBuf, cmd.Execute())
		})
	}
}

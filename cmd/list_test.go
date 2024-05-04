package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/antham/wo/workspace"
	"github.com/stretchr/testify/assert"
)

func TestNewListCmd(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T) workspaceManager
		test  func(*testing.T, *bytes.Buffer, *bytes.Buffer, error)
	}
	scenarios := []scenario{
		{
			"An error occurred when listing workspace",
			func(t *testing.T) workspaceManager {
				w := newMockWorkspaceManager(t)
				w.Mock.On("List").Return([]workspace.Workspace{}, errors.New("an error occurred"))
				return w
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.Error(t, err)
			},
		},
		{
			"No workspaces defined",
			func(t *testing.T) workspaceManager {
				w := newMockWorkspaceManager(t)
				w.Mock.On("List").Return([]workspace.Workspace{}, nil)
				return w
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "no workspaces defined")
			},
		},
		{
			"Listing workspaces",
			func(t *testing.T) workspaceManager {
				w := newMockWorkspaceManager(t)
				w.Mock.On("List").Return(
					[]workspace.Workspace{
						{
							Name: "api",
						},
						{
							Name: "db",
						},
					}, nil)
				return w
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.NoError(t, err)
				assert.Equal(t, outBuf.String(), `Workspaces
---
* api
* db
`)
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
			cmd := newListCmd(w)
			cmd.SetErr(errBuf)
			cmd.SetOut(outBuf)
			err := cmd.Execute()
			s.test(t, outBuf, errBuf, err)
		})
	}
}

package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/antham/wo/workspace"
	"github.com/stretchr/testify/assert"
)

func TestNewShowCmd(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T) (workspaceManager, []string)
		test  func(*testing.T, *bytes.Buffer, *bytes.Buffer, error)
	}
	scenarios := []scenario{
		{
			"An error occurred when showing a workspace",
			func(t *testing.T) (workspaceManager, []string) {
				w := newMockWorkspaceManager(t)
				args := []string{"api"}
				w.Mock.On("Get", args[0]).Return(workspace.Workspace{}, errors.New("an error occurred"))
				return w, args
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.Error(t, err)
			},
		},
		{
			"Showing a workspace with no functions and no envs",
			func(t *testing.T) (workspaceManager, []string) {
				w := newMockWorkspaceManager(t)
				args := []string{"api"}
				w.Mock.On("Get", args[0]).Return(workspace.Workspace{Name: args[0]}, nil)
				return w, args
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.NoError(t, err)
				assert.Equal(t, outBuf.String(), `Workspace api
---
Functions

No functions defined
---
Envs

No envs defined
`)
			},
		},
		{
			"Showing a workspace with functions and envs defined",
			func(t *testing.T) (workspaceManager, []string) {
				w := newMockWorkspaceManager(t)
				args := []string{"api"}
				w.Mock.On("Get", args[0]).Return(
					workspace.Workspace{
						Name: args[0],
						Functions: []workspace.Function{
							{
								Name:        "start",
								Description: "Start a server",
							},
							{
								Name:        "db-run",
								Description: "Start a db",
							},
							{
								Name: "stop",
							},
						},
						Envs: []string{
							"default",
							"dev",
							"prod",
						},
					}, nil)
				return w, args
			},
			func(t *testing.T, outBuf *bytes.Buffer, errBuf *bytes.Buffer, err error) {
				assert.NoError(t, err)
				assert.Equal(t, outBuf.String(), `Workspace api
---
Functions

* start : Start a server
* db-run : Start a db
* stop
---
Envs

* default
* dev
* prod
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
			w, args := s.setup(t)
			cmd := newShowCmd(w, newMockCompletionManager(t))
			cmd.SetArgs(args)
			cmd.SetErr(errBuf)
			cmd.SetOut(outBuf)
			err := cmd.Execute()
			s.test(t, outBuf, errBuf, err)
		})
	}
}

package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAliasCmd(t *testing.T) {
	type scenario struct {
		name  string
		args  []string
		setup func(*testing.T) workspaceManager
		test  func(*testing.T, *bytes.Buffer, *bytes.Buffer, error)
	}
	scenarios := []scenario{
		{
			"An error occurred when getting aliases",
			[]string{},
			func(t *testing.T) workspaceManager {
				w := newMockWorkspaceManager(t)
				w.Mock.On("BuildAliases", "").Return([]string{}, errors.New("an error occurred"))
				return w
			},
			func(t *testing.T, stdout *bytes.Buffer, stderr *bytes.Buffer, err error) {
				assert.Error(t, err)
			},
		},
		{
			"Listing aliases successfully",
			[]string{},
			func(t *testing.T) workspaceManager {
				w := newMockWorkspaceManager(t)
				w.Mock.On("BuildAliases", "").
					Return(
						[]string{
							`alias c_front="cd /tmp/front"`,
							`alias c_test="cd /tmp/test"`,
						},
						nil,
					)
				return w
			},
			func(t *testing.T, stdout *bytes.Buffer, stderr *bytes.Buffer, err error) {
				assert.NoError(t, err)
				assert.Equal(t,
					`alias c_front="cd /tmp/front"
alias c_test="cd /tmp/test"
`,
					stdout.String(),
				)
			},
		},
		{
			"Listing aliases with a prefix",
			[]string{"-p", "test_"},
			func(t *testing.T) workspaceManager {
				w := newMockWorkspaceManager(t)
				w.Mock.On("BuildAliases", "test_").
					Return(
						[]string{},
						nil,
					)
				return w
			},
			func(t *testing.T, stdout *bytes.Buffer, stderr *bytes.Buffer, err error) {
				assert.NoError(t, err)
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
			cmd := newAliasCmd(w)
			cmd.SetArgs(s.args)
			cmd.SetErr(errBuf)
			cmd.SetOut(outBuf)
			s.test(t, outBuf, errBuf, cmd.Execute())
		})
	}
}

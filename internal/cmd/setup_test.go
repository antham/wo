package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSetupCmd(t *testing.T) {
	type scenario struct {
		name  string
		args  []string
		setup func(*testing.T) workspaceManager
		test  func(*testing.T, *bytes.Buffer, *bytes.Buffer, error)
	}
	scenarios := []scenario{
		{
			"An error occurred when calling the command with a wrong first argument",
			[]string{"fis"},
			func(t *testing.T) workspaceManager {
				return newMockWorkspaceManager(t)
			},
			func(t *testing.T, stdout *bytes.Buffer, stderr *bytes.Buffer, err error) {
				assert.Error(t, err)
			},
		},
		{
			"An error occurred when getting aliases",
			[]string{"fish"},
			func(t *testing.T) workspaceManager {
				w := newMockWorkspaceManager(t)
				w.Mock.On("BuildAliases", "c_").Return([]string{}, errors.New("an error occurred"))
				return w
			},
			func(t *testing.T, stdout *bytes.Buffer, stderr *bytes.Buffer, err error) {
				assert.Error(t, err)
			},
		},
		{
			"We get the autocompletion for fish and aliases",
			[]string{"fish"},
			func(t *testing.T) workspaceManager {
				w := newMockWorkspaceManager(t)
				w.Mock.On("BuildAliases", "c_").
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
				assert.Contains(t,
					stdout.String(),
					`alias c_front="cd /tmp/front"
alias c_test="cd /tmp/test"`,
				)
				assert.Contains(t,
					stdout.String(),
					`function __wo_perform_completion`,
				)
			},
		},
		{
			"We get the autocompletion for bash and aliases",
			[]string{"bash"},
			func(t *testing.T) workspaceManager {
				w := newMockWorkspaceManager(t)
				w.Mock.On("BuildAliases", "c_").
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
				assert.Contains(t,
					stdout.String(),
					`alias c_front="cd /tmp/front"
alias c_test="cd /tmp/test"`,
				)
				assert.Contains(t,
					stdout.String(),
					`__wo_init_completion()`,
				)
			},
		},
		{
			"We get the autocompletion for zsh and aliases",
			[]string{"bash"},
			func(t *testing.T) workspaceManager {
				w := newMockWorkspaceManager(t)
				w.Mock.On("BuildAliases", "c_").
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
				assert.Contains(t,
					stdout.String(),
					`alias c_front="cd /tmp/front"
alias c_test="cd /tmp/test"`,
				)
				assert.Contains(t,
					stdout.String(),
					`__wo_init_completion()`,
				)
			},
		},
		{
			"We get the aliases only for sh",
			[]string{"sh"},
			func(t *testing.T) workspaceManager {
				w := newMockWorkspaceManager(t)
				w.Mock.On("BuildAliases", "c_").
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
			[]string{"fish", "-p", "test_"},
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
			cmd := newSetupCmd(w)
			cmd.SetArgs(s.args)
			cmd.SetErr(errBuf)
			cmd.SetOut(outBuf)
			s.test(t, outBuf, errBuf, cmd.Execute())
		})
	}
}

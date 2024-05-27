package completion

import (
	"errors"
	"testing"

	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCompletionProcess(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T) (Completion, string, []string)
		test  func(*testing.T, []string)
	}
	scenarios := []scenario{
		{
			"Args number greater than the number of decorators defined",
			func(t *testing.T) (Completion, string, []string) {
				w := newMockWorkspaceManager(t)
				c := New(w, []Decorator{})
				return c, "", []string{"test"}
			},
			func(t *testing.T, completions []string) {
				assert.Len(t, completions, 0)
			},
		},
		{
			"Return an error when calling the decorator",
			func(t *testing.T) (Completion, string, []string) {
				w := newMockWorkspaceManager(t)
				w.Mock.On("List").Return([]workspace.Workspace{}, errors.New("an error occurred"))
				c := New(w, []Decorator{FindWorkspaces})
				return c, "a", []string{}
			},
			func(t *testing.T, completions []string) {
				assert.Len(t, completions, 0)
			},
		},
		{
			"Define one decorator",
			func(t *testing.T) (Completion, string, []string) {
				w := newMockWorkspaceManager(t)
				w.Mock.On("List").Return(
					[]workspace.Workspace{
						{Name: "a"},
					}, nil)
				c := New(w, []Decorator{FindWorkspaces})
				return c, "a", []string{}
			},
			func(t *testing.T, completions []string) {
				assert.Len(t, completions, 1)
				assert.Equal(t, []string{"a"}, completions)
			},
		},
		{
			"Define two decorators",
			func(t *testing.T) (Completion, string, []string) {
				w := newMockWorkspaceManager(t)
				w.Mock.On("Get", "test").Return(
					workspace.Workspace{
						Envs: []string{
							"a",
						},
					}, nil)
				c := New(w, []Decorator{FindWorkspaces, FindEnvs})
				return c, "a", []string{"test"}
			},
			func(t *testing.T, completions []string) {
				assert.Len(t, completions, 1)
				assert.Equal(t, []string{"a"}, completions)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			c, toComplete, args := s.setup(t)
			completion, _ := c.Process(nil, args, toComplete)
			s.test(t, completion)
		})
	}
}

func TestFindWorkspaces(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T) (workspaceManager, string, []string)
		test  func(*testing.T, []string, cobra.ShellCompDirective, error)
	}
	scenarios := []scenario{
		{
			"An error occurred when listing workspaces",
			func(t *testing.T) (workspaceManager, string, []string) {
				w := newMockWorkspaceManager(t)
				w.Mock.On("List").Return([]workspace.Workspace{}, errors.New("an error occurred"))
				return w, "", []string{}
			},
			func(t *testing.T, completion []string, compMode cobra.ShellCompDirective, err error) {
				assert.Error(t, err)
				assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, compMode)
			},
		},
		{
			"Returns workspace matching the provided prefix",
			func(t *testing.T) (workspaceManager, string, []string) {
				w := newMockWorkspaceManager(t)
				w.Mock.On("List").Return(
					[]workspace.Workspace{
						{Name: "a"},
						{Name: "b"},
						{Name: "c"},
						{Name: "d"},
						{Name: "da"},
						{Name: "daa"},
						{Name: "dab"},
						{Name: "dac"},
						{Name: "db"},
						{Name: "dc"},
						{Name: "e"},
					}, nil)
				return w, "da", []string{}
			},
			func(t *testing.T, completion []string, compMode cobra.ShellCompDirective, err error) {
				assert.NoError(t, err)
				assert.Equal(t, []string{"da", "daa", "dab", "dac"}, completion)
				assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, compMode)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			workspaceManager, toComplete, args := s.setup(t)
			completion, compMode, err := FindWorkspaces(workspaceManager, toComplete, args...)
			s.test(t, completion, compMode, err)
		})
	}
}

func TestFindEnvs(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T) (workspaceManager, string, []string)
		test  func(*testing.T, []string, cobra.ShellCompDirective, error)
	}
	scenarios := []scenario{
		{
			"An error occurred when getting envs",
			func(t *testing.T) (workspaceManager, string, []string) {
				w := newMockWorkspaceManager(t)
				w.Mock.On("Get", "test").Return(workspace.Workspace{}, errors.New("an error occurred"))
				return w, "", []string{"test"}
			},
			func(t *testing.T, completion []string, compMode cobra.ShellCompDirective, err error) {
				assert.Error(t, err)
				assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, compMode)
			},
		},
		{
			"Returns envs matching the provided prefix",
			func(t *testing.T) (workspaceManager, string, []string) {
				w := newMockWorkspaceManager(t)
				w.Mock.On("Get", "test").Return(
					workspace.Workspace{
						Envs: []string{
							"a",
							"b",
							"c",
							"d",
							"da",
							"daa",
							"dab",
							"dac",
							"db",
							"dc",
							"e",
						},
					}, nil)
				return w, "da", []string{"test"}
			},
			func(t *testing.T, completion []string, compMode cobra.ShellCompDirective, err error) {
				assert.NoError(t, err)
				assert.Equal(t, []string{"da", "daa", "dab", "dac"}, completion)
				assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, compMode)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			workspaceManager, toComplete, args := s.setup(t)
			completion, compMode, err := FindEnvs(workspaceManager, toComplete, args...)
			s.test(t, completion, compMode, err)
		})
	}
}

func TestFindFunctions(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T) (workspaceManager, string, []string)
		test  func(*testing.T, []string, cobra.ShellCompDirective, error)
	}
	scenarios := []scenario{
		{
			"An error occurred when getting functions",
			func(t *testing.T) (workspaceManager, string, []string) {
				w := newMockWorkspaceManager(t)
				w.Mock.On("Get", "test").Return(workspace.Workspace{}, errors.New("an error occurred"))
				return w, "", []string{"test"}
			},
			func(t *testing.T, completion []string, compMode cobra.ShellCompDirective, err error) {
				assert.Error(t, err)
				assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, compMode)
			},
		},
		{
			"Returns functions matching the provided prefix",
			func(t *testing.T) (workspaceManager, string, []string) {
				w := newMockWorkspaceManager(t)
				w.Mock.On("Get", "test").Return(
					workspace.Workspace{
						Functions: []workspace.Function{
							{Name: "a"},
							{Name: "b"},
							{Name: "c"},
							{Name: "d"},
							{Name: "da"},
							{Name: "daa", Description: "function daa"},
							{Name: "dab"},
							{Name: "dac", Description: "function dac"},
							{Name: "db"},
							{Name: "dc"},
							{Name: "e"},
						},
					}, nil)
				return w, "da", []string{"test"}
			},
			func(t *testing.T, completion []string, compMode cobra.ShellCompDirective, err error) {
				assert.NoError(t, err)
				assert.Equal(t, []string{"da", "daa\tfunction daa", "dab", "dac\tfunction dac"}, completion)
				assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, compMode)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			workspaceManager, toComplete, args := s.setup(t)
			completion, compMode, err := FindFunctions(workspaceManager, toComplete, args...)
			s.test(t, completion, compMode, err)
		})
	}
}

func TestFindDir(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T) (workspaceManager, string, []string)
		test  func(*testing.T, []string, cobra.ShellCompDirective, error)
	}
	scenarios := []scenario{
		{
			"Get completion dirs directive",
			func(t *testing.T) (workspaceManager, string, []string) {
				return newMockWorkspaceManager(t), "", []string{"test"}
			},
			func(t *testing.T, completion []string, compMode cobra.ShellCompDirective, err error) {
				assert.NoError(t, err)
				assert.Empty(t, completion)
				assert.Equal(t, cobra.ShellCompDirectiveFilterDirs, compMode)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			workspaceManager, toComplete, args := s.setup(t)
			completion, compMode, err := FindDirs(workspaceManager, toComplete, args...)
			s.test(t, completion, compMode, err)
		})
	}
}

func TestNoOp(t *testing.T) {
	type scenario struct {
		name  string
		setup func(*testing.T) (workspaceManager, string, []string)
		test  func(*testing.T, []string, cobra.ShellCompDirective, error)
	}
	scenarios := []scenario{
		{
			"Get no file directive",
			func(t *testing.T) (workspaceManager, string, []string) {
				return newMockWorkspaceManager(t), "", []string{"test"}
			},
			func(t *testing.T, completion []string, compMode cobra.ShellCompDirective, err error) {
				assert.NoError(t, err)
				assert.Empty(t, completion)
				assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, compMode)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			workspaceManager, toComplete, args := s.setup(t)
			completion, compMode, err := NoOp(workspaceManager, toComplete, args...)
			s.test(t, completion, compMode, err)
		})
	}
}

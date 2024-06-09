package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVersionCmd(t *testing.T) {
	out := &bytes.Buffer{}
	cmd := newVersionCmd()
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetOut(out)
	assert.NoError(t, cmd.Execute())
	assert.Equal(t, "dev\n", out.String())
}

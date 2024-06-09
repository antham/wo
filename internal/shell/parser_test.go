package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	fs, err := Parse("bash", []byte(`
# This is a function to run
function f1 {
	echo e;
}

function function_test {
	echo e;
}
`))
	assert.NoError(t, err)
	assert.Len(t, fs, 2)
	assert.Equal(t, fs, []Function{
		{Name: "f1", Description: "This is a function to run"},
		{Name: "function_test", Description: ""},
	})
	fs, err = Parse("sh", []byte(`
# This is a function to run
function f1 {
	echo e;
}

function function_test {
	echo e;
}
`))
	assert.NoError(t, err)
	assert.Len(t, fs, 2)
	assert.Equal(t, fs, []Function{
		{Name: "f1", Description: "This is a function to run"},
		{Name: "function_test", Description: ""},
	})
	fs, err = Parse("fish", []byte(`
function f1 -d "This is a function to run"
	echo e
end

function f2
	echo e
end
`))
	assert.NoError(t, err)
	assert.Len(t, fs, 2)
	assert.Equal(t, fs, []Function{
		{Name: "f1", Description: "This is a function to run"},
		{Name: "f2", Description: ""},
	})
	fs, err = Parse("zsh", []byte(`
# This is a function to run
function f1 {
	echo e
}
function f2 () {
	echo e
}
`))
	assert.NoError(t, err)
	assert.Len(t, fs, 2)
	assert.Equal(t, fs, []Function{
		{Name: "f1", Description: "This is a function to run"},
		{Name: "f2", Description: ""},
	})
}

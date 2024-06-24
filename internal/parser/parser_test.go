package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	fs := Parse("bash", []byte(`
# This is a function to run
f1() {
	echo e;
}

function_test() {
	echo e;
}
`))
	assert.Len(t, fs, 2)
	assert.Equal(t, fs, []Function{
		{Name: "f1", Description: "This is a function to run"},
		{Name: "function_test", Description: ""},
	})
	fs = Parse("sh", []byte(`
# This is a function to run
f1 () {
	echo e;
}

function_test() {
	echo e;
}
`))
	assert.Len(t, fs, 2)
	assert.Equal(t, fs, []Function{
		{Name: "f1", Description: "This is a function to run"},
		{Name: "function_test", Description: ""},
	})
	fs = Parse("fish", []byte(`
function f1 -d "This is a function to run"
	echo e
end

function f2
	echo e
end
`))
	assert.Len(t, fs, 2)
	assert.Equal(t, fs, []Function{
		{Name: "f1", Description: "This is a function to run"},
		{Name: "f2", Description: ""},
	})
	fs = Parse("zsh", []byte(`
# This is a function to run
f1 () {
	echo e
}
f2 () {
	echo e
}
`))
	assert.Len(t, fs, 2)
	assert.Equal(t, fs, []Function{
		{Name: "f1", Description: "This is a function to run"},
		{Name: "f2", Description: ""},
	})
}

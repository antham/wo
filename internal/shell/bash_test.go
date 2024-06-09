package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBashParser(t *testing.T) {
	bashParser := newBashParser()
	data, err := bashParser.parse([]byte(`
f1() {
    echo e;
}

f2() { echo e;}

function f3{
    echo e;
}

function f4{ echo e;}

# This is a description comment
f5() {
    echo e;
}

# This is a description comment
f6() { echo e;}

# This is a description comment
function f7{
    echo e;
}

# This is a description comment
function f8{ echo e;}
`))
	functions := data
	assert.NoError(t, err)
	assert.Len(t, functions, 8)
	assert.Equal(t, []Function{
		{Name: "f1", Description: ""},
		{Name: "f2", Description: ""},
		{Name: "f3", Description: ""},
		{Name: "f4", Description: ""},
		{Name: "f5", Description: "This is a description comment"},
		{Name: "f6", Description: "This is a description comment"},
		{Name: "f7", Description: "This is a description comment"},
		{Name: "f8", Description: "This is a description comment"},
	}, functions)
}

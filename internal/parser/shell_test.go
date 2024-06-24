package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellParser(t *testing.T) {
	shellParser := newShellParser()
	functions := shellParser.parse([]byte(`
f1() {
    echo e;
}

f2() { echo e;}

# This is a description comment
f3() {
    echo e;
}

# This is a description comment
f4() { echo e;}
`))
	assert.Len(t, functions, 4)
	assert.Equal(t, []Function{
		{Name: "f1", Description: ""},
		{Name: "f2", Description: ""},
		{Name: "f3", Description: "This is a description comment"},
		{Name: "f4", Description: "This is a description comment"},
	}, functions)
}

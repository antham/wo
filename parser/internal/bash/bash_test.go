package bash

import (
	"testing"

	"github.com/antham/wo/parser"
	"github.com/stretchr/testify/assert"
)

func TestBashParser(t *testing.T) {
	bashParser := newBashParser()
	data, err := bashParser.Parse([]byte(`
# This is a description comment
function f1() {
	echo e;
}



function f2() { echo e;}

f3() {
   echo e;
}

f4() { echo e;}


# This is any comment

# On several lines
f5  () { echo e;}

another_little_func () {}

function_test () {}
`))
	functions := data.([]parser.Function)
	assert.NoError(t, err)
	assert.Len(t, functions, 7)
	assert.Equal(t, []parser.Function{
		{Name: "f1", Description: "This is a description comment"},
		{Name: "f2"},
		{Name: "f3"},
		{Name: "f4"},
		{Name: "f5"},
		{Name: "another_little_func"},
		{Name: "function_test"},
	}, functions)
}

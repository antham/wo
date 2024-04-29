package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFishParser(t *testing.T) {
	fishParser := newFishParser()
	data, err := fishParser.parse([]byte(`
function f1 -d "f1 description comment"
	echo e
end

function f2
	echo e
end

function f3 --description "f3 description comment"
	echo e
end

function f4 --description "f4 description comment";echo e; end

function f5 -d "f5 description comment";echo e; end

function f6;echo e; end

function f7 -d "function to do something";echo e; end
`))
	functions := data
	assert.NoError(t, err)
	assert.Len(t, functions, 7)
	assert.Equal(t, []Function{
		{Name: "f1", Description: "f1 description comment"},
		{Name: "f2"},
		{Name: "f3", Description: "f3 description comment"},
		{Name: "f4", Description: "f4 description comment"},
		{Name: "f5", Description: "f5 description comment"},
		{Name: "f6"},
		{Name: "f7", Description: "function to do something"},
	}, functions)
}

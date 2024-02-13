package fish

import (
	"testing"

	"github.com/antham/wo/parser"
	"github.com/stretchr/testify/assert"
)

func TestFishParser(t *testing.T) {
	fishParser := newFishParser()
	data, err := fishParser.Parse([]byte(`
function f1 -d "f1 description command"
	echo e
end

function f2
	echo e
end

function f3 --description "f3 description command"
	echo e
end

function f4 --description "f4 description command";echo e; end

function f5 -d "f5 description command";echo e; end

function f6;echo e; end
`))
	functions := data.([]parser.Function)
	assert.NoError(t, err)
	assert.Len(t, functions, 9)
	assert.Equal(t, []parser.Function{
		{Name: "f1", Description: "f1 description comment"},
		{Name: "f2"},
		{Name: "f3", Description: "f3 description comment"},
		{Name: "f4", Description: "f4 description comment"},
		{Name: "f5", Description: "f5 description comment"},
	}, functions)
}

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRubyParser(t *testing.T) {
	rubyParser := newRubyParser()
	functions := rubyParser.parse([]byte(`
def f1
  puts "Hello Ruby"
end

def f2(name)
  puts "hello #\{name}"
end

# f3 description comment
def f3
  puts "Hello Ruby"
end

# f4 description comment
def f4(name)
  puts "hello #\{name}"
end
`))
	assert.Len(t, functions, 4)
	assert.Equal(t, []Function{
		{Name: "f1"},
		{Name: "f2"},
		{Name: "f3", Description: "f3 description comment"},
		{Name: "f4", Description: "f4 description comment"},
	}, functions)
}

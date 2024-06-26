package parser

import (
	"regexp"
	"strings"
)

type rubyParser struct{}

func newRubyParser() *rubyParser {
	return &rubyParser{}
}

func (r *rubyParser) parse(content []byte) []Function {
	functions := []Function{}
	re := regexp.MustCompile(`(?m)(?:#\s*(?P<description>.+)\n)?(?:def\s+(?P<function>[^\n|(]+))`)
	for _, match := range re.FindAllSubmatch(content, -1) {
		function := Function{}
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				switch name {
				case "function":
					if len(match[i]) == 0 {
						continue
					}
					function.Name = strings.TrimSpace(string(match[i]))
				case "description":
					function.Description = strings.TrimSpace(string(match[i]))
				}
			}
		}
		functions = append(functions, function)
	}
	return functions
}

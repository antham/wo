package parser

import (
	"regexp"
	"strings"
)

type shellParser struct{}

func newShellParser() *shellParser {
	return &shellParser{}
}

func (shellParser *shellParser) parse(content []byte) []Function {
	functions := []Function{}
	r := regexp.MustCompile(`(?m)(?:#\s*(?P<description>.+)\n)?(?:(?P<function>.*)\s*\(\))(?:\s*|\n)?{`)
	for _, match := range r.FindAllSubmatch(content, -1) {
		function := Function{}
		for i, name := range r.SubexpNames() {
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

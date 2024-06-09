package shell

import (
	"regexp"
	"strings"
)

type bashParser struct {
}

func newBashParser() *bashParser {
	return &bashParser{}
}

func (bashParser *bashParser) parse(content []byte) ([]Function, error) {
	functions := []Function{}
	r := regexp.MustCompile(`(?m)(?:#\s*(?P<description>.+)\n)?(?:function\s+?(?P<function_1>.+)|(?P<function_2>.*)\s*\(\))(?:\s*|\n)?{`)
	for _, match := range r.FindAllSubmatch(content, -1) {
		function := Function{}
		for i, name := range r.SubexpNames() {
			if i != 0 && name != "" {
				switch name {
				case "function_1", "function_2":
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
	return functions, nil
}

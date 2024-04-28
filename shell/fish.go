package shell

import (
	"regexp"
)

type fishParser struct {
}

func newFishParser() *fishParser {
	return &fishParser{}
}

func (fishParser *fishParser) parse(content []byte) ([]Function, error) {
	fs := []Function{}
	r := regexp.MustCompile(`function\s+(?P<function>[^ |;|\n]+)(?:\s+--?d(?:escription)?\s+(?:"|')(?P<definition>[^(?:'|")]+)(?:"|'))?`).FindAllStringSubmatch(string(content), -1)
	for _, l := range r {
		f := Function{Name: l[1]}
		if len(l) == 3 {
			f.Description = l[2]
		}
		fs = append(fs, f)
	}
	return fs, nil
}

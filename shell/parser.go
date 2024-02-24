package shell

import (
	"strings"

	"github.com/bzick/tokenizer"
)

const (
	TokenCurlyOpen tokenizer.TokenKey = iota + 1
	TokenCurlyClose
	TokenParenOpen
	TokenParenClose
	TokenComment
	TokenDescriptionShort
	TokenDescriptionLong
	TokenQuote
	TokenFunction
	TokenSemiColon
)

type Shell string

const (
	Zsh  Shell = "zsh"
	Bash Shell = "bash"
	Fish Shell = "fish"
)

type Function struct {
	Name        string
	Description string
}

func Parse(shell string, content []byte) ([]Function, error) {
	var data interface{}
	var err error
	switch shell {
	case string(Zsh):
		p := newZshParser()
		data, err = p.parse(content)
		if err != nil {
			return []Function{}, nil
		}
	case string(Bash):
		p := newBashParser()
		data, err = p.parse(content)
		if err != nil {
			return []Function{}, nil
		}
	case string(Fish):
		p := newFishParser()
		data, err = p.parse(content)
		if err != nil {
			return []Function{}, nil
		}
	}
	fs := []Function{}
	for _, f := range data.([]Function) {
		f.Description = strings.ReplaceAll(f.Description, "function ", "function")
		fs = append(fs, f)
	}
	return fs, nil
}

func createDescription(descriptionTokens []*tokenizer.Token) string {
	description := ""
	for i, d := range descriptionTokens {
		if i > 0 && descriptionTokens[i-1].Offset()+1 == d.Offset() || i == 0 {
			description += d.ValueString()
		} else {
			description += " " + d.ValueString()
		}
	}
	return description
}

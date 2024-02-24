package shell

import (
	"strings"

	"github.com/bzick/tokenizer"
)

const (
	tokenCurlyOpen tokenizer.TokenKey = iota + 1
	tokenCurlyClose
	tokenParenOpen
	tokenParenClose
	tokenComment
	tokenDescriptionShort
	tokenDescriptionLong
	tokenQuote
	tokenFunction
	tokenSemiColon
)

type shellStr string

const (
	zsh  shellStr = "zsh"
	bash shellStr = "bash"
	fish shellStr = "fish"
)

type Function struct {
	Name        string
	Description string
}

func Parse(shell string, content []byte) ([]Function, error) {
	var data interface{}
	var err error
	switch shell {
	case string(zsh):
		p := newZshParser()
		data, err = p.parse(content)
		if err != nil {
			return []Function{}, nil
		}
	case string(bash):
		p := newBashParser()
		data, err = p.parse(content)
		if err != nil {
			return []Function{}, nil
		}
	case string(fish):
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

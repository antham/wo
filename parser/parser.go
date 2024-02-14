package parser

import "github.com/bzick/tokenizer"

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
	Bash       = "bash"
	Fish       = "fish"
)

type Function struct {
	Name        string
	Description string
}

func Parse(shell string, content []byte) ([]Function, error) {
	switch shell {
	case string(Zsh):
		p := newZshParser()
		fs, err := p.parse(content)
		return fs.([]Function), err
	case string(Bash):
		p := newBashParser()
		fs, err := p.parse(content)
		return fs.([]Function), err
	case string(Fish):
		p := newFishParser()
		fs, err := p.parse(content)
		return fs.([]Function), err
	}
	return []Function{}, nil
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

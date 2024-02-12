package bash

import (
	"strings"

	"github.com/antham/wo/parser"
	"github.com/bzick/tokenizer"
)

const (
	TokenCurlyOpen tokenizer.TokenKey = iota + 1
	TokenCurlyClose
	TokenParenOpen
	TokenParenClose
	TokenFunction
	TokenComment
)

type bashParser struct {
	tokenizer *tokenizer.Tokenizer
}

func newBashParser() *bashParser {
	bashParser := &bashParser{}
	bashParser.tokenizer = tokenizer.New()
	bashParser.tokenizer.
		DefineTokens(TokenCurlyOpen, []string{"{"}).
		DefineTokens(TokenCurlyClose, []string{"}"}).
		DefineTokens(TokenParenOpen, []string{"("}).
		DefineTokens(TokenParenClose, []string{")"}).
		DefineTokens(TokenFunction, []string{"function "}).
		DefineTokens(TokenComment, []string{"#"})
	return bashParser
}

func (bashParser *bashParser) Parse(content []byte) (interface{}, error) {
	return bashParser.analyzer(bashParser.tokenizer.ParseBytes(content))
}

func (bashParser *bashParser) analyzer(stream *tokenizer.Stream) (interface{}, error) {
	functions := []parser.Function{}
	comments := map[int][]string{}
	for {
		if stream.CurrentToken().Key() == TokenComment {
			currentToken := stream.CurrentToken()
			stream.GoNext()
			for {
				if stream.CurrentToken().Line() != currentToken.Line() {
					break
				}
				comments[stream.CurrentToken().Line()] = append(comments[stream.CurrentToken().Line()], stream.CurrentToken().ValueString())
				stream.GoNext()
			}
		}
		if stream.IsNextSequence(TokenParenOpen, TokenParenClose) {
			currentToken := stream.CurrentToken()
			acc := ""
			for {
				if stream.CurrentToken().Line() != currentToken.Line() || stream.CurrentToken().Key() == TokenFunction {
					break
				}
				acc = stream.CurrentToken().ValueString() + acc
				stream.GoPrev()
			}
			stream.GoTo(currentToken.ID())
			functions = append(functions, parser.Function{Name: acc, Description: strings.Join(comments[stream.CurrentToken().Line()-1], " ")})
		}
		stream.GoNext()
		if stream.CurrentToken() == nil {
			break
		}
	}
	return functions, nil
}

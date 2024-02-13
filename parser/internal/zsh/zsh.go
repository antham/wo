package zsh

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

type zshParser struct {
	tokenizer *tokenizer.Tokenizer
}

func newZshParser() *zshParser {
	zshParser := &zshParser{}
	zshParser.tokenizer = tokenizer.New()
	zshParser.tokenizer.
		DefineTokens(TokenCurlyOpen, []string{"{"}).
		DefineTokens(TokenCurlyClose, []string{"}"}).
		DefineTokens(TokenParenOpen, []string{"("}).
		DefineTokens(TokenParenClose, []string{")"}).
		DefineTokens(TokenFunction, []string{"function "}).
		DefineTokens(TokenComment, []string{"#"})
	return zshParser
}

func (zshParser *zshParser) Parse(content []byte) (interface{}, error) {
	return zshParser.analyzer(zshParser.tokenizer.ParseBytes(content))
}

func (zshParser *zshParser) analyzer(stream *tokenizer.Stream) (interface{}, error) {
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
		if stream.CurrentToken().Key() == TokenCurlyOpen && stream.PrevToken().Key() != TokenParenClose {
			currentToken := stream.CurrentToken()
			stream.GoPrev()
			isFunction := false
			for {
				if stream.CurrentToken().Line() != currentToken.Line() {
					break
				}
				if stream.CurrentToken().Key() == TokenFunction {
					stream.GoNext()
					isFunction = true
					break
				}
				stream.GoPrev()
			}
			if isFunction {
				acc := ""
				for {
					if stream.CurrentToken().Key() == TokenCurlyOpen {
						break
					}
					acc = acc + stream.CurrentToken().ValueString()
					stream.GoNext()
				}
				stream.GoTo(currentToken.ID())
				functions = append(functions, parser.Function{Name: acc, Description: strings.Join(comments[stream.CurrentToken().Line()-1], " ")})
			}
			stream.GoTo(currentToken.ID())
		}

		stream.GoNext()
		if stream.CurrentToken() == nil {
			break
		}
	}
	return functions, nil
}

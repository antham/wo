package parser

import (
	"github.com/bzick/tokenizer"
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

func (bashParser *bashParser) parse(content []byte) (interface{}, error) {
	return bashParser.analyzer(bashParser.tokenizer.ParseBytes(content))
}

func (bashParser *bashParser) analyzer(stream *tokenizer.Stream) (interface{}, error) {
	functions := []Function{}
	comments := map[int][]*tokenizer.Token{}
	for {
		if stream.CurrentToken().Key() == TokenComment {
			currentToken := stream.CurrentToken()
			stream.GoNext()
			for {
				if stream.CurrentToken().Line() != currentToken.Line() {
					break
				}
				comments[stream.CurrentToken().Line()] = append(comments[stream.CurrentToken().Line()], stream.CurrentToken())
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
			functions = append(functions, Function{Name: acc, Description: createDescription(comments[stream.CurrentToken().Line()-1])})
		}
		stream.GoNext()
		if stream.CurrentToken().Key() == 0 {
			break
		}
	}
	return functions, nil
}

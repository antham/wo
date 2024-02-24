package shell

import (
	"github.com/bzick/tokenizer"
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

func (zshParser *zshParser) parse(content []byte) (interface{}, error) {
	return zshParser.analyzer(zshParser.tokenizer.ParseBytes(content))
}

func (zshParser *zshParser) analyzer(stream *tokenizer.Stream) (interface{}, error) {
	functions := []Function{}
	comments := map[int][]*tokenizer.Token{}
	for {
		if stream.CurrentToken() == nil || stream.CurrentToken().Key() == 0 {
			break
		}
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
				functions = append(functions, Function{Name: acc, Description: createDescription(comments[stream.CurrentToken().Line()-1])})
			}
			stream.GoTo(currentToken.ID())
		}
		stream.GoNext()
	}
	return functions, nil
}

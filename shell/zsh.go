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
		DefineTokens(tokenCurlyOpen, []string{"{"}).
		DefineTokens(tokenCurlyClose, []string{"}"}).
		DefineTokens(tokenParenOpen, []string{"("}).
		DefineTokens(tokenParenClose, []string{")"}).
		DefineTokens(tokenFunction, []string{"function "}).
		DefineTokens(tokenComment, []string{"#"})
	return zshParser
}

func (zshParser *zshParser) parse(content []byte) ([]Function, error) {
	return zshParser.analyzer(zshParser.tokenizer.ParseBytes(content))
}

func (zshParser *zshParser) analyzer(stream *tokenizer.Stream) ([]Function, error) {
	functions := []Function{}
	comments := map[int][]*tokenizer.Token{}
	for {
		if stream.CurrentToken() == nil || stream.CurrentToken().Key() == 0 {
			break
		}
		if stream.CurrentToken().Key() == tokenComment {
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
		if stream.IsNextSequence(tokenParenOpen, tokenParenClose) {
			currentToken := stream.CurrentToken()
			acc := ""
			for {
				if stream.CurrentToken().Line() != currentToken.Line() || stream.CurrentToken().Key() == tokenFunction {
					break
				}
				acc = stream.CurrentToken().ValueString() + acc
				stream.GoPrev()
			}
			stream.GoTo(currentToken.ID())
			functions = append(functions, Function{Name: acc, Description: createDescription(comments[stream.CurrentToken().Line()-1])})
		}
		if stream.CurrentToken().Key() == tokenCurlyOpen && stream.PrevToken().Key() != tokenParenClose {
			currentToken := stream.CurrentToken()
			stream.GoPrev()
			isFunction := false
			for {
				if stream.CurrentToken().Line() != currentToken.Line() {
					break
				}
				if stream.CurrentToken().Key() == tokenFunction {
					stream.GoNext()
					isFunction = true
					break
				}
				stream.GoPrev()
			}
			if isFunction {
				acc := ""
				for {
					if stream.CurrentToken().Key() == tokenCurlyOpen {
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

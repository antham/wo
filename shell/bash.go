package shell

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
		DefineTokens(tokenCurlyOpen, []string{"{"}).
		DefineTokens(tokenCurlyClose, []string{"}"}).
		DefineTokens(tokenParenOpen, []string{"("}).
		DefineTokens(tokenParenClose, []string{")"}).
		DefineTokens(tokenFunction, []string{"function "}).
		DefineTokens(tokenComment, []string{"#"})
	return bashParser
}

func (bashParser *bashParser) parse(content []byte) (interface{}, error) {
	return bashParser.analyzer(bashParser.tokenizer.ParseBytes(content))
}

func (bashParser *bashParser) analyzer(stream *tokenizer.Stream) (interface{}, error) {
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
		stream.GoNext()
	}
	return functions, nil
}

package parser

import (
	"github.com/bzick/tokenizer"
)

type fishParser struct {
	tokenizer *tokenizer.Tokenizer
}

func newFishParser() *fishParser {
	fishParser := &fishParser{}
	fishParser.tokenizer = tokenizer.New()
	fishParser.tokenizer.
		DefineTokens(TokenDescriptionShort, []string{"-d"}).
		DefineTokens(TokenDescriptionShort, []string{"--description"}).
		DefineTokens(TokenQuote, []string{`"`}).
		DefineTokens(TokenFunction, []string{"function "}).
		DefineTokens(TokenSemiColon, []string{";"})
	return fishParser
}

func (fishParser *fishParser) parse(content []byte) (interface{}, error) {
	return fishParser.analyzer(fishParser.tokenizer.ParseBytes(content))
}

func (fishParser *fishParser) analyzer(stream *tokenizer.Stream) (interface{}, error) {
	functions := []Function{}
	for {
		if stream.CurrentToken() == nil || stream.CurrentToken().Key() == 0 {
			break
		}
		if stream.CurrentToken().Key() == TokenFunction {
			tokenFunction := stream.CurrentToken()
			function := ""
			for {
				stream.GoNext()
				if stream.CurrentToken().Key() == TokenDescriptionLong || stream.CurrentToken().Key() == TokenDescriptionShort {
					break
				}
				if stream.CurrentToken().Line() != tokenFunction.Line() || stream.CurrentToken().Key() == TokenSemiColon {
					break
				}
				function += stream.CurrentToken().ValueString()
			}
			descriptionTokens := []*tokenizer.Token{}
			if stream.CurrentToken().Key() == TokenDescriptionLong || stream.CurrentToken().Key() == TokenDescriptionShort {
				stream.GoNext()
				if stream.CurrentToken().Key() == TokenQuote {
					for {
						stream.GoNext()
						if stream.CurrentToken().Key() == TokenQuote {
							break
						}
						descriptionTokens = append(descriptionTokens, stream.CurrentToken())
					}
				}

			}
			functions = append(functions, Function{Name: function, Description: createDescription(descriptionTokens)})
		}
		stream.GoNext()
	}
	return functions, nil
}

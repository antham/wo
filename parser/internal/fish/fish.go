package fish

import (
	"strings"

	"github.com/antham/wo/parser"
	"github.com/bzick/tokenizer"
)

const (
	TokenDescriptionShort tokenizer.TokenKey = iota + 1
	TokenDescriptionLong
	TokenQuote
	TokenFunction
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
		DefineTokens(TokenFunction, []string{"function "})
	return fishParser
}

func (fishParser *fishParser) Parse(content []byte) (interface{}, error) {
	return fishParser.analyzer(fishParser.tokenizer.ParseBytes(content))
}

func (fishParser *fishParser) analyzer(stream *tokenizer.Stream) (interface{}, error) {
	functions := []parser.Function{}
	for {
		if stream.CurrentToken().Key() == TokenFunction {
			stream.GoNext()
			function := stream.CurrentToken().ValueString()
			description := []string{""}
			stream.GoNext()
			if stream.CurrentToken().Key() == TokenDescriptionLong || stream.CurrentToken().Key() == TokenDescriptionShort {
				stream.GoNext()
				if stream.CurrentToken().Key() == TokenQuote {
					for {
						stream.GoNext()
						if stream.CurrentToken().Key() == TokenQuote {
							break
						}
						description = append(description, stream.CurrentToken().ValueString())
					}
				}

			}
			functions = append(functions, parser.Function{Name: function, Description: strings.Join(description, " ")})
		}

		stream.GoNext()
		if stream.CurrentToken() == nil {
			break
		}
	}
	return functions, nil
}

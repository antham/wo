package parser

type appStr string

const (
	zsh  appStr = "zsh"
	bash appStr = "bash"
	fish appStr = "fish"
	sh   appStr = "sh"
	ruby appStr = "rb"
)

type Function struct {
	Name        string
	Description string
}

func Parse(app string, content []byte) []Function {
	switch app {
	case string(bash), string(sh), string(zsh):
		return newShellParser().parse(content)
	case string(fish):
		return newFishParser().parse(content)
	case string(ruby):
		return newRubyParser().parse(content)
	}
	return []Function{}
}

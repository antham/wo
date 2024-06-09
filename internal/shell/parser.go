package shell

type shellStr string

const (
	zsh  shellStr = "zsh"
	bash shellStr = "bash"
	fish shellStr = "fish"
	sh   shellStr = "sh"
)

type Function struct {
	Name        string
	Description string
}

func Parse(shell string, content []byte) []Function {
	switch shell {
	case string(bash), string(sh), string(zsh):
		return newShellParser().parse(content)
	case string(fish):
		return newFishParser().parse(content)
	}
	return []Function{}
}

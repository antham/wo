package workspace

type Commander interface {
	command(string, ...string) error
}

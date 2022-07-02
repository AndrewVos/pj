package actions

type Action interface {
	Apply() error
}

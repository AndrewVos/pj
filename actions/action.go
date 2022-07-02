package actions

type Action interface {
	Flag() string
	AddActionDescription() string
	Apply() error
}

var All []Action

func RegisterAction(action Action) {
	All = append(All, action)
}

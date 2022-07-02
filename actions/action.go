package actions

type Action interface {
	Flag() string
	AddActionDescription() string
	Apply(modulePath string) error
	Completions(fieldName string) ([]string, error)
}

var All []Action

func RegisterAction(action Action) {
	All = append(All, action)
}

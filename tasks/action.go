package tasks

type Task interface {
	Flag() string
	AddTaskDescription() string
	Apply(modulePath string) error
	Completions(fieldName string) ([]string, error)
}

var All []Task

func RegisterTask(task Task) {
	All = append(All, task)
}

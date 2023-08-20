package components

type Task struct {
	Status      Status
	title       string
	description string
	id          string
}

func NewTask(status Status, title, description string, id string) Task {
	return Task{Status: status, title: title, description: description, id: id}
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

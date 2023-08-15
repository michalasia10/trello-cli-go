package components

type Task struct {
	Status      Status
	title       string
	description string
}

func NewTask(status Status, title, description string) Task {
	return Task{Status: status, title: title, description: description}
}

func (t *Task) Next() {
	if t.Status == Done {
		t.Status = Todo
	} else {
		t.Status++
	}
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

package components

type Status int

func (s Status) GetNext() Status {
	if s == Done {
		return Todo
	}
	return s + 1
}

func (s Status) GetPrev() Status {
	if s == Todo {
		return Done
	}
	return s - 1
}

const Margin = 4

const (
	Todo Status = iota
	InProgress
	Done
)

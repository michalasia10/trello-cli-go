package components

type Status int

func (s Status) GetNext(maxColumns int) Status {
	newS := s + 1
	if newS >= Status(maxColumns) {
		return 0
	}
	return newS

}

func (s Status) GetPrev(maxColumns int) Status {
	newS := s - 1
	if newS < Status(0) {
		return Status(maxColumns) - 1
	}
	return newS
}

const Margin = 4

const (
	Todo Status = iota
	InProgress
	Done
)

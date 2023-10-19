package fd

type State interface {
	isState()
}

type state int

func (state) isState() {}

const (
	StateUnknown = state(iota)
	StateFileInfo
	StateRead
)

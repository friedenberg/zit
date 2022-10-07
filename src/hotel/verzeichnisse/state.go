package verzeichnisse

type State int

const (
	StateUnread = State(iota)
	StateRead
	StateChanged
)


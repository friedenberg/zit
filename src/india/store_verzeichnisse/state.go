package store_verzeichnisse

type State int

const (
	StateUnread = State(iota)
	StateReadJustHeader
	StateRead
	StateChanged
)

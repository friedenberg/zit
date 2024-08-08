package external_state

//go:generate stringer -type=State
type State int

const (
	Unknown = State(iota)
	Tracked
	Untracked
	Recognized
	Conflicted
	Error
)

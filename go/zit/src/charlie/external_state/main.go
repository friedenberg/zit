package external_state

import "fmt"

type State int

const (
	Unknown = State(iota)
	Tracked
	Untracked
	Recognized
	Deleted
	WouldDelete
)

func (s State) String() string {
	switch s {
	case Tracked:
		return "tracked"

	case Untracked:
		return "untracked"

	case Recognized:
		return "recognized"

	case Deleted:
		return "deleted"

	case WouldDelete:
		return "would delete"

	default:
		return fmt.Sprintf("Unknown(%d)", s)
	}
}

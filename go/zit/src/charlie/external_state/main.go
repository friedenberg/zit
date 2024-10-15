package external_state

import "fmt"

type State int

const (
	Unknown = State(iota)
	Tracked
	Untracked
	Recognized
	Conflicted
	Parent
	Deleted
	WouldDelete
	Error
)

func (s State) String() string {
	switch s {
	case Tracked:
		return "tracked"

	case Untracked:
		return "untracked"

	case Recognized:
		return "recognized"

	case Conflicted:
		return "conflict"

	case Parent:
		return "parent"

	case Deleted:
		return "deleted"

	case WouldDelete:
		return "would delete"

	case Error:
		return "error"

	default:
		return fmt.Sprintf("Unknown(%d)", s)
	}
}

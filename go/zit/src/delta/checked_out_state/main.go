package checked_out_state

import (
	"fmt"

	"code.linenisgreat.com/zit/src/delta/string_format_writer"
)

type State int

const (
	StateUnknown = State(iota)
	StateEmpty
	StateJustCheckedOut
	StateJustCheckedOutButDifferent
	StateExistsAndSame
	StateExistsAndDifferent
	StateUntracked
	StateRecognized
	StateConflicted
	StateError
)

func (s State) String() string {
	switch s {
	case StateJustCheckedOut:
		return string_format_writer.StringCheckedOut

	case StateExistsAndSame:
		return string_format_writer.StringSame

	case StateJustCheckedOutButDifferent:
		fallthrough
	case StateExistsAndDifferent:
		return string_format_writer.StringChanged

	case StateUntracked:
		return string_format_writer.StringUntracked

	case StateRecognized:
		return string_format_writer.StringRecognized

	case StateConflicted:
		return string_format_writer.StringConflicted

	case StateError:
		return "error"

	default:
		return fmt.Sprintf("unknown: %#v", s)
	}
}

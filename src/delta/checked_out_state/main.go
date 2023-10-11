package checked_out_state

import (
	"fmt"

	"github.com/friedenberg/zit/src/charlie/string_format_writer"
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

	default:
		return fmt.Sprintf("unknown: %#v", s)
	}
}

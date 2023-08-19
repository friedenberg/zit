package checked_out_state

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/string_writer_format"
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
)

func (s State) String() string {
	switch s {
	case StateJustCheckedOut:
		return string_writer_format.StringCheckedOut

	case StateExistsAndSame:
		return string_writer_format.StringSame

	case StateJustCheckedOutButDifferent:
		fallthrough
	case StateExistsAndDifferent:
		return string_writer_format.StringChanged

	case StateUntracked:
		return string_writer_format.StringUntracked

	case StateRecognized:
		return string_writer_format.StringRecognized

	default:
		return fmt.Sprintf("unknown: %#v", s)
	}
}

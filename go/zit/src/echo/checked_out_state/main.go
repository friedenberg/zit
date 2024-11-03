package checked_out_state

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
)

type State int

const (
	Unknown = State(iota)
	JustCheckedOut
	JustCheckedOutButDifferent
	ExistsAndSame
	ExistsAndDifferent
	Untracked
	Recognized
	Conflicted
	BlobMissing
	Error
)

func (s State) String() string {
	switch s {
	case JustCheckedOut:
		return string_format_writer.StringCheckedOut

	case ExistsAndSame:
		return string_format_writer.StringSame

	case JustCheckedOutButDifferent:
		fallthrough
	case ExistsAndDifferent:
		return string_format_writer.StringChanged

	case Untracked:
		return string_format_writer.StringUntracked

	case Recognized:
		return string_format_writer.StringRecognized

	case Conflicted:
		return string_format_writer.StringConflicted

	case Error:
		return "error"

	case BlobMissing:
		return string_format_writer.StringBlobMissing

	default:
		return fmt.Sprintf("unknown: %#v", s)
	}
}

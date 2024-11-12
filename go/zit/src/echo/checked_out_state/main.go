package checked_out_state

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
)

// TODO define this state much more clearly, as it's currently overloaded and
// abused
type State int

/*

State       | Internal | External | Equality
------------|----------|----------|---------
empty       | none     | none     | invalid
transacted  | some     | none     | invalid
untracked   | none     | some     | invalid
checked out | some     | some     | valid

*/

const (
	Unknown        = State(iota)
	JustCheckedOut // UI
	CheckedOut     // UI
	ExistsAndSame  // Internal v External
	Changed        // Internal v External
	Untracked      // Internal v External
	Recognized     // Internal v External
	Conflicted     // Internal v External
)

func (s State) String() string {
	switch s {
	case JustCheckedOut:
		return string_format_writer.StringCheckedOut

	case ExistsAndSame:
		return string_format_writer.StringSame

	case Changed:
		return string_format_writer.StringChanged

	case Untracked:
		return string_format_writer.StringUntracked

	case Recognized:
		return string_format_writer.StringRecognized

	case Conflicted:
		return string_format_writer.StringConflicted

	default:
		return fmt.Sprintf("unknown: %#v", s)
	}
}

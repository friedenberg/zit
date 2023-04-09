package objekte

import (
	"fmt"
)

type CheckedOutState int

const (
	CheckedOutStateUnknown = CheckedOutState(iota)
	CheckedOutStateEmpty
	CheckedOutStateJustCheckedOut
	CheckedOutStateJustCheckedOutButSame
	CheckedOutStateExistsAndSame
	CheckedOutStateExistsAndDifferent
	CheckedOutStateUntracked
	CheckedOutStateRecognized
)

func (s CheckedOutState) String() string {
	switch s {
	case CheckedOutStateJustCheckedOutButSame:
		return "checked out"

	case CheckedOutStateExistsAndSame:
		return "same"

	case CheckedOutStateExistsAndDifferent:
		return "changed"

	case CheckedOutStateUntracked:
		return "untracked"

	case CheckedOutStateRecognized:
		return "recognized"

	default:
		return fmt.Sprintf("unknown: %#v", s)
	}
}

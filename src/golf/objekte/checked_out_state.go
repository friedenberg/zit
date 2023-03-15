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

	default:
		return fmt.Sprintf("unknown: %#v", s)
	}
}

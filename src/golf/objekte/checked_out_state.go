package objekte

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
	case CheckedOutStateExistsAndSame, CheckedOutStateJustCheckedOutButSame:
		return "same"

	case CheckedOutStateExistsAndDifferent:
		return "changed"

	case CheckedOutStateUntracked:
		return "untracked"

	default:
		return "unknown"
	}
}

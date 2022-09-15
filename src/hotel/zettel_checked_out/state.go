package zettel_checked_out

type State int

const (
	StateUnknown = State(iota)
	StateEmpty
	StateJustCheckedOut
	StateJustCheckedOutButSame
	StateExistsAndSame
	StateExistsAndDifferent
	StateUntracked
)

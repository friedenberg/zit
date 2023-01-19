package zettel_checked_out

type State int

const (
	StateNotCheckedOut = State(iota)
	StateEmpty
	StateJustCheckedOut
	StateJustCheckedOutButSame
	StateExistsAndSame
	StateExistsAndDifferent
	StateUntracked
)

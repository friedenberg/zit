package zettel_checked_out

type State int

const (
	StateUnknown = State(iota)
  StateEmpty
	StateExistsAndSame
	StateExistsAndDifferent
	StateUntracked
)

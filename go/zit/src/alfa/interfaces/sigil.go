package interfaces

// Used primarily for limiting / restricting queries.
type Sigil interface {
	IncludesHistory
	IncludesLatest
	IncludesExternal
	IncludesHidden
}

type SigilGetter interface {
	GetSigil() Sigil
}

type IncludesHistory interface {
	IncludesHistory() bool
}

type IncludesLatest interface {
	IncludesLatest() bool
}

type IncludesExternal interface {
	IncludesExternal() bool
}

type IncludesHidden interface {
	IncludesHidden() bool
}

package schnittstellen

type Sigil interface {
	IncludesHistory
	IncludesSchwanzen
	IncludesExternal
	IncludesHidden
}

type SigilGetter interface {
	GetSigil() Sigil
}

type IncludesHistory interface {
	IncludesHistory() bool
}

type IncludesSchwanzen interface {
	IncludesSchwanzen() bool
}

type IncludesExternal interface {
	IncludesExternal() bool
}

type IncludesHidden interface {
	IncludesHidden() bool
}

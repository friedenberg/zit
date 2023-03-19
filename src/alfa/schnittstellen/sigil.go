package schnittstellen

type Sigil interface {
	IncludesHistory
	IncludesSchwanzen
	IncludesCwd
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

type IncludesCwd interface {
	IncludesCwd() bool
}

type IncludesHidden interface {
	IncludesHidden() bool
}

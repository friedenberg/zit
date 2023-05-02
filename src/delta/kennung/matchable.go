package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type EtikettenGetter interface {
	GetEtiketten() schnittstellen.Set[Etikett]
}

type EtikettenExpandedGetter interface {
	GetEtikettenExpanded() schnittstellen.Set[Etikett]
}

type TypGetter interface {
	GetTyp() Typ
}

type IdLikeGetter interface {
	GetIdLike() Kennung
}

type Matchable interface {
	schnittstellen.ValueLike
	schnittstellen.Stored

	EtikettenGetter
	EtikettenExpandedGetter
	TypGetter
	IdLikeGetter
}

type MatchableGetter interface {
	GetMatchable() Matchable
}

type MatchableAdder interface {
	AddMatchable(Matchable) error
}

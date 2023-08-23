package matcher

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type EtikettenGetter interface {
	GetEtiketten() kennung.EtikettSet
}

type TypGetter interface {
	GetTyp() kennung.Typ
}

type IdLikeGetter interface {
	GetKennungLike() kennung.Kennung
	GetKennungLikePtr() kennung.KennungPtr
}

type Matchable interface {
	schnittstellen.ValueLike
	schnittstellen.Stored

	EtikettenGetter
	TypGetter
	IdLikeGetter
}

type MatchableGetter interface {
	GetMatchable() Matchable
}

type MatchableAdder interface {
	AddMatchable(Matchable) error
}

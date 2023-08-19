package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type EtikettenGetter interface {
	GetEtiketten() EtikettSet
}

type TypGetter interface {
	GetTyp() Typ
}

type IdLikeGetter interface {
	GetKennungLike() Kennung
	GetKennungLikePtr() KennungPtr
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

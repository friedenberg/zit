package kennung

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type EtikettenGetter interface {
	GetEtiketten() *schnittstellen.Set[Etikett]
}

type etikettenGetter struct{}

func (_ etikettenGetter) GetEtiketten() *schnittstellen.Set[Etikett] {
	return nil
}

type TypGetter interface {
	GetTyp() *Typ
}

type typGetter struct{}

func (_ typGetter) GetTyp() *Typ {
	return nil
}

type IdLikeGetter interface {
	GetIdLike() IdLike
}

type Matchable interface {
	EtikettenGetter
	TypGetter
	IdLikeGetter
}

type MatchableGetter interface {
	GetMatchable() Matchable
}

type matchable struct {
	EtikettenGetter
	TypGetter
	IdLikeGetter
}

func MakeMatchable(
	etiketten EtikettenGetter,
	typ TypGetter,
	id IdLikeGetter,
) Matchable {
	if etiketten == nil {
		etiketten = etikettenGetter{}
	}

	if typ == nil {
		typ = typGetter{}
	}

	if id == nil {
		panic("nil IdLikeGetter for Matchable")
	}

	return matchable{
		EtikettenGetter: etiketten,
		TypGetter:       typ,
		IdLikeGetter:    id,
	}
}

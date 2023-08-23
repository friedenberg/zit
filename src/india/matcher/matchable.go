package matcher

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type IdLikeGetter interface {
	GetKennungLike() kennung.Kennung
	GetKennungLikePtr() kennung.KennungPtr
}

type Matchable interface {
	schnittstellen.ValueLike
	schnittstellen.Stored

	sku.SkuLikePtr
}

type MatchableGetter interface {
	GetMatchable() Matchable
}

type MatchableAdder interface {
	AddMatchable(Matchable) error
}

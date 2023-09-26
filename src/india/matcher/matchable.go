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

	*sku.Transacted
}

type MatchableGetter interface {
	GetMatchable() *sku.Transacted
}

type MatchableAdder interface {
	AddMatchable(*sku.Transacted) error
}

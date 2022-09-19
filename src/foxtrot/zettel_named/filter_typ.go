package zettel_named

import (
	"github.com/friedenberg/zit/src/charlie/typ"
)

type FilterTyp typ.Typ

func (f FilterTyp) IncludeNamedZettel(z Zettel) (ok bool) {
	ok = typ.Typ(f).Includes(z.Stored.Zettel.Typ.Etikett)
	return
}

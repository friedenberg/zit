package zettel_named

import (
	"github.com/friedenberg/zit/src/echo/typ"
)

type FilterTyp typ.Kennung

func (f FilterTyp) IncludeNamedZettel(z Zettel) (ok bool) {
	ok = typ.Kennung(f).Includes(z.Stored.Objekte.Typ)
	return
}

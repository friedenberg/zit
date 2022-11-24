package zettel_named

import "github.com/friedenberg/zit/src/charlie/kennung"

type FilterTyp kennung.Typ

func (f FilterTyp) IncludeNamedZettel(z Zettel) (ok bool) {
	ok = kennung.Typ(f).Includes(z.Stored.Objekte.Typ)
	return
}

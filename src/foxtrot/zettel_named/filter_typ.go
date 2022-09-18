package zettel_named

import "github.com/friedenberg/zit/src/charlie/typ"

type FilterTyp typ.Typ

func (f FilterTyp) IncludeNamedZettel(z Zettel) bool {
	return typ.Typ(f).Contains(z.Stored.Zettel.Typ)
}

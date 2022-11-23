package zettel_named

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_stored"
)

type Zettel struct {
	Stored  zettel_stored.Stored
	Kennung hinweis.Hinweis
}

func (a *Zettel) Equals(b *Zettel) bool {
	if !a.Stored.Equals(&b.Stored) {
		errors.Print("stored")
		return false
	}

	if !a.Kennung.Equals(b.Kennung) {
		errors.Print("hinweis")
		return false
	}

	return true
}

func (zn *Zettel) Reset() {
	zn.Kennung = hinweis.Hinweis{}
	zn.Stored.Reset()
}

// func (zn *Zettel) LineFormat(k konfig.Konfig) zettel_line.Format {
//   f := zettel_line.MakeFromKonfig(k).Builder().

// 	zi := p.MakeZettelish().
// 		Id(p.Hinweis(zn.Hinweis)).
// 		Rev(p.Sha(zn.Stored.Sha)).
// 		TypString(zn.Stored.Zettel.Typ.String()).
// 		Bez(p.Bezeichnung(zn.Stored.Zettel))
// }

package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
)

func init() {
	errors.TodoP0("add hidden")
}

type Verzeichnisse struct {
	wasPopulated bool
	Etiketten    kennung_index.EtikettenVerzeichnisse
	Typ          kennung_index.TypVerzeichnisse
	// Hidden bool
}

func (z *Verzeichnisse) ResetWithObjekte(z1 Objekte) {
	z.wasPopulated = true

	z.Etiketten.ResetWithEtikettSet(z1.Etiketten)
	z.Typ.ResetWithTyp(z1.Typ)
}

func (z *Verzeichnisse) Reset() {
	z.wasPopulated = false

	z.Etiketten.Reset()
	z.Typ.Reset()
}

func (z *Verzeichnisse) ResetWith(z1 Verzeichnisse) {
	z.wasPopulated = true

	z.Etiketten.ResetWith(z1.Etiketten)
	z.Typ.ResetWith(z1.Typ)
}

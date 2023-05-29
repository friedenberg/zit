package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

func init() {
	errors.TodoP4("add hidden")
}

type Verzeichnisse struct {
	wasPopulated bool
	Typ          kennung_index.TypVerzeichnisse
	// Hidden bool
}

func (z *Verzeichnisse) ResetWithObjekteMetadateiGetter(
	z1 Objekte,
	mg metadatei.Getter,
) {
	z.wasPopulated = true

	m := mg.GetMetadatei()
	z.Typ.ResetWithTyp(m.GetTyp())
}

func (z *Verzeichnisse) Reset() {
	z.wasPopulated = false

	z.Typ.Reset()
}

func (z *Verzeichnisse) ResetWith(z1 Verzeichnisse) {
	z.wasPopulated = true

	z.Typ.ResetWith(z1.Typ)
}

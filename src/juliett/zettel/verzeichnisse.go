package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
)

func init() {
	errors.TodoP4("add hidden")
}

type Verzeichnisse struct {
	wasPopulated bool
	Etiketten    kennung_index.EtikettenVerzeichnisse
	Typ          kennung_index.TypVerzeichnisse
	// Hidden bool
}

func (z Verzeichnisse) GetEtiketten() schnittstellen.Set[kennung.Etikett] {
	return z.Etiketten.GetEtiketten()
}

func (z Verzeichnisse) GetEtikettenExpanded() schnittstellen.Set[kennung.Etikett] {
	return z.Etiketten.GetEtikettenExpanded()
}

func (z *Verzeichnisse) ResetWithObjekte(z1 Objekte) {
	z.wasPopulated = true

	z.Etiketten.ResetWithEtikettSet(z1.Metadatei.Etiketten)
	z.Typ.ResetWithTyp(z1.GetMetadatei().GetTyp())
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

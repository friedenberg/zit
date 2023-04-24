package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
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

func (z *Verzeichnisse) ResetWithObjekteMetadateiGetter(
	z1 Objekte,
	mg metadatei.Getter,
) {
	z.wasPopulated = true

	m := mg.GetMetadatei()
	z.Etiketten.ResetWithEtikettSet(m.Etiketten)
	z.Typ.ResetWithTyp(m.GetTyp())
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

package metadatei

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
)

type Metadatei struct {
	Bezeichnung bezeichnung.Bezeichnung
	Etiketten   kennung.EtikettSet
	Typ         kennung.Typ
}

func (z Metadatei) GetTyp() kennung.Typ {
	return z.Typ
}

func (z Metadatei) GetEtiketten() schnittstellen.Set[kennung.Etikett] {
	return z.Etiketten.ImmutableClone()
}

func (z Metadatei) Equals(z1 Metadatei) bool {
	if !z.Typ.Equals(z1.Typ) {
		return false
	}

	if !z.Bezeichnung.Equals(z1.Bezeichnung) {
		return false
	}

	if !z.Etiketten.Equals(z1.Etiketten) {
		return false
	}

	return true
}

func (z *Metadatei) Reset() {
	z.Typ = kennung.Typ{}
	z.Bezeichnung.Reset()
	z.Etiketten = kennung.MakeEtikettSet()
}

func (z *Metadatei) ResetWith(z1 Metadatei) {
	z.Typ = z1.Typ
	z.Bezeichnung = z1.Bezeichnung
	z.Etiketten = z1.Etiketten.ImmutableClone()
}

func (z Metadatei) Description() (d string) {
	d = z.Bezeichnung.String()

	if strings.TrimSpace(d) == "" {
		d = collections.StringCommaSeparated[kennung.Etikett](z.Etiketten)
	}

	return
}

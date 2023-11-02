package metadatei

import "github.com/friedenberg/zit/src/echo/kennung"

var Resetter resetter

type resetter struct{}

func (resetter) Reset(z *Metadatei) {
	z.AkteSha.Reset()
	z.Bezeichnung.Reset()
	z.Etiketten = kennung.MakeEtikettSet()
	z.Verzeichnisse.Reset()
	z.Typ = kennung.Typ{}
	// z.Gattung = gattung.Unknown
	z.Tai.Reset()
}

func (resetter) ResetWith(z *Metadatei, z1 Metadatei) {
	z.AkteSha = z1.AkteSha
	z.Bezeichnung = z1.Bezeichnung

	if z1.Etiketten == nil {
		z.Etiketten = kennung.MakeEtikettSet()
	} else {
		z.Etiketten = z1.Etiketten.CloneSetPtrLike()
	}

	z.Verzeichnisse.ResetWith(&z1.Verzeichnisse)

	z.Typ = z1.Typ
	z.Tai = z1.Tai
}

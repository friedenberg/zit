package metadatei

import (
	"github.com/friedenberg/zit/src/echo/kennung"
)

var Resetter resetter

type resetter struct{}

func (resetter) Reset(z *Metadatei) {
	z.AkteSha.Reset()
	z.Bezeichnung.Reset()
	z.SetEtiketten(nil)
	ResetterVerzeichnisse.Reset(&z.Verzeichnisse)
	z.Typ = kennung.Typ{}
	// z.Gattung = gattung.Unknown
	z.Tai.Reset()
}

func (resetter) ResetWith(z *Metadatei, z1 Metadatei) {
	z.AkteSha = z1.AkteSha
	z.Bezeichnung = z1.Bezeichnung

	z.SetEtiketten(z1.Etiketten)

	ResetterVerzeichnisse.ResetWith(&z.Verzeichnisse, &z1.Verzeichnisse)

	z.Typ = z1.Typ
	z.Tai = z1.Tai
}

var ResetterVerzeichnisse resetterVerzeichnisse

type resetterVerzeichnisse struct{}

func (resetterVerzeichnisse) Reset(v *Verzeichnisse) {
	v.Archiviert.Reset()
	v.SetExpandedEtiketten(nil)
	v.SetImplicitEtiketten(nil)
	v.Mutter.Reset()
	v.Sha.Reset()
}

func (resetterVerzeichnisse) ResetWith(a, b *Verzeichnisse) {
	a.Archiviert.ResetWith(b.Archiviert)
	a.SetExpandedEtiketten(b.GetExpandedEtiketten())
	a.SetImplicitEtiketten(b.GetImplicitEtiketten())
	a.Mutter = b.Mutter
	a.Sha = b.Sha
}

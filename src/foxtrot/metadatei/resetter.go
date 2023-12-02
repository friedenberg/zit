package metadatei

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/kennung"
)

var Resetter resetter

type resetter struct{}

func (resetter) Reset(z *Metadatei) {
	z.AkteSha.Reset()
	z.Bezeichnung.Reset()
	z.Comments = z.Comments[:0]
	z.SetEtiketten(nil)
	ResetterVerzeichnisse.Reset(&z.Verzeichnisse)
	z.Typ = kennung.Typ{}
	// z.Gattung = gattung.Unknown
	z.Tai.Reset()
}

func (r resetter) ResetWith(a *Metadatei, b Metadatei) {
	r.ResetWithPtr(a, &b)
}

func (resetter) ResetWithPtr(a *Metadatei, b *Metadatei) {
	errors.PanicIfError(a.AkteSha.SetShaLike(b.AkteSha))
	a.Bezeichnung = b.Bezeichnung
	a.Comments = a.Comments[:0]
	a.Comments = append(a.Comments, b.Comments...)

	a.SetEtiketten(b.Etiketten)

	ResetterVerzeichnisse.ResetWith(&a.Verzeichnisse, &b.Verzeichnisse)

	a.Typ = b.Typ
	a.Tai = b.Tai
}

var ResetterVerzeichnisse resetterVerzeichnisse

type resetterVerzeichnisse struct{}

func (resetterVerzeichnisse) Reset(a *Verzeichnisse) {
	a.Archiviert.Reset()
	a.SetExpandedEtiketten(nil)
	a.SetImplicitEtiketten(nil)
	a.Mutter.Reset()
	a.Sha.Reset()
}

func (resetterVerzeichnisse) ResetWith(a, b *Verzeichnisse) {
	a.Archiviert.ResetWith(b.Archiviert)
	a.SetExpandedEtiketten(b.GetExpandedEtiketten())
	a.SetImplicitEtiketten(b.GetImplicitEtiketten())
	a.Mutter = b.Mutter
	a.Sha = b.Sha
}

package metadatei

import (
	"slices"

	"code.linenisgreat.com/zit/src/echo/kennung"
)

var Resetter resetter

type resetter struct{}

func (resetter) Reset(z *Metadatei) {
	z.Bezeichnung.Reset()
	z.Comments = z.Comments[:0]
	z.ResetEtiketten()
	ResetterVerzeichnisse.Reset(&z.Verzeichnisse)
	z.Typ = kennung.Typ{}
	// z.Gattung = gattung.Unknown
	z.Tai.Reset()
	z.Shas.Reset()
}

func (resetter) ResetWith(a *Metadatei, b *Metadatei) {
	a.Bezeichnung = b.Bezeichnung
	a.Comments = a.Comments[:0]
	a.Comments = append(a.Comments, b.Comments...)

	a.SetEtiketten(b.Etiketten)

	ResetterVerzeichnisse.ResetWith(&a.Verzeichnisse, &b.Verzeichnisse)

	a.Typ = b.Typ
	a.Tai = b.Tai

	a.Shas.ResetWith(&b.Shas)
}

var ResetterVerzeichnisse resetterVerzeichnisse

type resetterVerzeichnisse struct{}

func (resetterVerzeichnisse) Reset(a *Verzeichnisse) {
	a.Etiketten.Reset()
	a.Archiviert.Reset()
	a.SetExpandedEtiketten(nil)
	a.SetImplicitEtiketten(nil)
	a.QueryPath.Reset()
}

func (resetterVerzeichnisse) ResetWith(a, b *Verzeichnisse) {
	a.Etiketten.ResetWith(&b.Etiketten)
	a.Archiviert.ResetWith(b.Archiviert)
	a.SetExpandedEtiketten(b.GetExpandedEtiketten())
	a.SetImplicitEtiketten(b.GetImplicitEtiketten())
	a.QueryPath.Reset()
	a.QueryPath = slices.Grow(a.QueryPath, b.QueryPath.Len())
	copy(a.QueryPath, b.QueryPath)
}

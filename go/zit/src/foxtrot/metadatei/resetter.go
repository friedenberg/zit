package metadatei

import (
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/etiketten_path"
)

var Resetter resetter

type resetter struct{}

func (resetter) Reset(z *Metadatei) {
	z.Bezeichnung.Reset()
	z.Comments = z.Comments[:0]
	z.SetEtiketten(nil)
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
	a.Etiketten = a.Etiketten[:0]
	a.Archiviert.Reset()
	a.SetExpandedEtiketten(nil)
	a.SetImplicitEtiketten(nil)
}

func (resetterVerzeichnisse) ResetWith(a, b *Verzeichnisse) {
	a.Etiketten = make([]*etiketten_path.Path, len(b.Etiketten))
	copy(a.Etiketten, b.Etiketten)
	a.Archiviert.ResetWith(b.Archiviert)
	a.SetExpandedEtiketten(b.GetExpandedEtiketten())
	a.SetImplicitEtiketten(b.GetImplicitEtiketten())
}

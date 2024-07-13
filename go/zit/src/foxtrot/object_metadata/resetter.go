package object_metadata

import (
	"slices"

	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

var Resetter resetter

type resetter struct{}

func (resetter) Reset(z *Metadata) {
	z.Description.Reset()
	z.Comments = z.Comments[:0]
	z.ResetEtiketten()
	ResetterVerzeichnisse.Reset(&z.Cached)
	z.Type = ids.Type{}
	// z.Gattung = gattung.Unknown
	z.Tai.Reset()
	z.Shas.Reset()
}

func (resetter) ResetWith(a *Metadata, b *Metadata) {
	a.Description = b.Description
	a.Comments = a.Comments[:0]
	a.Comments = append(a.Comments, b.Comments...)

	a.SetEtiketten(b.Tags)

	ResetterVerzeichnisse.ResetWith(&a.Cached, &b.Cached)

	a.Type = b.Type
	a.Tai = b.Tai

	a.Shas.ResetWith(&b.Shas)
}

var ResetterVerzeichnisse resetterVerzeichnisse

type resetterVerzeichnisse struct{}

func (resetterVerzeichnisse) Reset(a *Verzeichnisse) {
	a.Etiketten.Reset()
	a.Schlummernd.Reset()
	a.SetExpandedEtiketten(nil)
	a.SetImplicitEtiketten(nil)
	a.QueryPath.Reset()
}

func (resetterVerzeichnisse) ResetWith(a, b *Verzeichnisse) {
	a.Etiketten.ResetWith(&b.Etiketten)
	a.Schlummernd.ResetWith(b.Schlummernd)
	a.SetExpandedEtiketten(b.GetExpandedEtiketten())
	a.SetImplicitEtiketten(b.GetImplicitEtiketten())
	a.QueryPath.Reset()
	a.QueryPath = slices.Grow(a.QueryPath, b.QueryPath.Len())
	copy(a.QueryPath, b.QueryPath)
}

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
	z.ResetTags()
	ResetterVerzeichnisse.Reset(&z.Cache)
	z.Type = ids.Type{}
	// z.Gattung = gattung.Unknown
	z.Tai.Reset()
	z.Shas.Reset()
}

func (resetter) ResetWith(a *Metadata, b *Metadata) {
	a.Description = b.Description
	a.Comments = a.Comments[:0]
	a.Comments = append(a.Comments, b.Comments...)

	a.SetTags(b.Tags)

	ResetterVerzeichnisse.ResetWith(&a.Cache, &b.Cache)

	a.Type = b.Type
	a.Tai = b.Tai

	a.Shas.ResetWith(&b.Shas)
}

var ResetterVerzeichnisse resetterVerzeichnisse

type resetterVerzeichnisse struct{}

func (resetterVerzeichnisse) Reset(a *Cache) {
	a.TagPaths.Reset()
	a.Dormant.Reset()
	a.SetExpandedTags(nil)
	a.SetImplicitTags(nil)
	a.QueryPath.Reset()
}

func (resetterVerzeichnisse) ResetWith(a, b *Cache) {
	a.TagPaths.ResetWith(&b.TagPaths)
	a.Dormant.ResetWith(b.Dormant)
	a.SetExpandedTags(b.GetExpandedTags())
	a.SetImplicitTags(b.GetImplicitTags())
	a.QueryPath.Reset()
	a.QueryPath = slices.Grow(a.QueryPath, b.QueryPath.Len())
	copy(a.QueryPath, b.QueryPath)
}

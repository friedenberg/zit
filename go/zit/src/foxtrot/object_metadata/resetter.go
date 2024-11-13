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
	ResetterCache.Reset(&z.Cache)
	z.Type = ids.Type{}
	z.Tai.Reset()
	z.Shas.Reset()
	z.Fields = z.Fields[:0]
}

func (resetter) ResetWithExceptFields(dst *Metadata, src *Metadata) {
	dst.Description = src.Description
	dst.Comments = dst.Comments[:0]
	dst.Comments = append(dst.Comments, src.Comments...)

	dst.SetTags(src.Tags)

	ResetterCache.ResetWith(&dst.Cache, &src.Cache)

	dst.Type = src.Type
	dst.Tai = src.Tai

	dst.Shas.ResetWith(&src.Shas)
}

func (r resetter) ResetWith(dst *Metadata, src *Metadata) {
	r.ResetWithExceptFields(dst, src)
	dst.Fields = dst.Fields[:0]
	dst.Fields = append(dst.Fields, src.Fields...)
}

var ResetterCache resetterCache

type resetterCache struct{}

func (resetterCache) Reset(a *Cache) {
	a.ParentTai.Reset()
	a.TagPaths.Reset()
	a.Dormant.Reset()
	a.SetExpandedTags(nil)
	a.SetImplicitTags(nil)
	a.QueryPath.Reset()
}

func (resetterCache) ResetWith(a, b *Cache) {
	a.ParentTai.ResetWith(b.ParentTai)
	a.TagPaths.ResetWith(&b.TagPaths)
	a.Dormant.ResetWith(b.Dormant)
	a.SetExpandedTags(b.GetExpandedTags())
	a.SetImplicitTags(b.GetImplicitTags())
	a.QueryPath.Reset()
	a.QueryPath = slices.Grow(a.QueryPath, b.QueryPath.Len())
	copy(a.QueryPath, b.QueryPath)
}

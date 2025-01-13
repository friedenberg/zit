package query

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/india/sku"
)

type ObjectId struct {
	Exact   bool
	Virtual bool
	Debug   bool

	*ids.ObjectId
}

func (k ObjectId) reduce(b *buildState) (err error) {
	if err = k.GetObjectId().Expand(b.builder.expanders); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO support exact
func (exp ObjectId) ContainsSku(tg sku.TransactedGetter) (ok bool) {
	sk := tg.GetSku()

	defer sk.Metadata.Cache.QueryPath.PushOnReturn(exp, &ok)

	skMe := sk.GetMetadata()

	method := ids.Contains

	if exp.Exact {
		method = ids.ContainsExactly
	}

	switch exp.GetGenre() {
	case genres.Tag:
		var idx int

		if exp.Exact {
			idx, ok = skMe.Cache.TagPaths.All.ContainsObjectIdTagExact(
				exp.GetObjectId(),
			)
		} else {
			idx, ok = skMe.Cache.TagPaths.All.ContainsObjectIdTag(
				exp.GetObjectId(),
			)
		}

		ui.Log().Print(exp, idx, ok, skMe.Cache.TagPaths.All, sk)

		if ok {
			ps := skMe.Cache.TagPaths.All[idx]
			sk.Metadata.Cache.QueryPath.PushOnReturn(ps.Parents, &ok)
			return
		}

		return

	case genres.Type:
		if method(skMe.GetType(), exp.GetObjectId()) {
			ok = true
			return
		}

		if e, isExternal := tg.(*sku.Transacted); isExternal {
			if method(e.ExternalType, exp.GetObjectId()) {
				ok = true
				return
			}
		}
	}

	idl := &sk.ObjectId

	if !method(idl, exp.GetObjectId()) {
		return
	}

	ok = true

	return
}

func (k ObjectId) String() string {
	var sb strings.Builder

	if k.Exact {
		sb.WriteRune('=')
	}

	if k.Virtual {
		sb.WriteRune('%')
	}

	sb.WriteString(ids.FormattedString(k.GetObjectId()))

	return sb.String()
}

package query

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type ObjectId struct {
	Exact    bool
	Virtual  bool
	Debug    bool
	External bool

	*ids.ObjectId
}

func (k ObjectId) Reduce(b *Builder) (err error) {
	if err = k.Expand(b.expanders); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO support exact
func (k ObjectId) ContainsSku(sk *sku.Transacted) (ok bool) {
	defer sk.Metadata.Cache.QueryPath.PushOnOk(k, &ok)

	me := sk.GetMetadata()

	switch k.GetGenre() {
	case genres.Tag:
		var idx int

		if k.Exact {
			idx, ok = me.Cache.TagPaths.All.ContainsObjectIdTagExact(
				k.ObjectId,
			)
		} else {
			idx, ok = me.Cache.TagPaths.All.ContainsObjectIdTag(
				k.ObjectId,
			)
		}

		ui.Log().Print(k, idx, ok, me.Cache.TagPaths.All, sk)

		if ok {
			// if k.Exact {
			// 	ewp := me.Verzeichnisse.Etiketten.All[idx]
			// 	ui.Debug().Print(ewp, sk)
			// }

			ps := me.Cache.TagPaths.All[idx]
			sk.Metadata.Cache.QueryPath.Push(ps.Parents)
			return
		}

		return

	case genres.Type:
		if ids.Contains(me.GetType(), k) {
			ok = true
			return
		}

		// case kennung.ShaLike:
		// 	if Sha(kt.GetSha()).ContainsMatchable(m) {
		// 		return true
		// 	}
	}

	idl := &sk.ObjectId

	if !ids.Contains(idl, k) {
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

	sb.WriteString(ids.FormattedString(k.ObjectId))

	return sb.String()
}

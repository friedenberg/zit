package query

import (
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Kennung struct {
	useEtikettenPaths bool

	Exact   bool
	Virtual bool
	Debug   bool
	FD      *fd.FD
	*kennung.Kennung2
}

func (k Kennung) Reduce(b *Builder) (err error) {
	if err = k.Expand(b.expanders); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO support exact
func (k Kennung) ContainsSku(sk *sku.Transacted) (ok bool) {
	defer sk.Metadatei.Verzeichnisse.QueryPath.PushOnOk(k, &ok)

	me := sk.GetMetadatei()
	switch k.GetGattung() {
	case gattung.Etikett:
		if k.useEtikettenPaths {
			kps := k.PartsStrings()

			var idx int
			idx, ok = me.Verzeichnisse.Etiketten.ContainsEtikett(kps.Right)

			if ok {
				ps := me.Verzeichnisse.Etiketten.All[idx]
				sk.Metadatei.Verzeichnisse.QueryPath.Push(ps.Parents)
				return
			}
		} else {
			s := k.String()

			if me.GetEtiketten().ContainsKey(s) {
				ok = true
				return
			}

			if me.Verzeichnisse.GetExpandedEtiketten().ContainsKey(s) {
				ok = true
				return
			}

			if me.Verzeichnisse.GetImplicitEtiketten().ContainsKey(s) {
				ok = true
				return
			}
		}

	case gattung.Typ:
		if kennung.Contains(me.GetTyp(), k) {
			ok = true
			return
		}

		// case kennung.ShaLike:
		// 	if Sha(kt.GetSha()).ContainsMatchable(m) {
		// 		return true
		// 	}
	}

	idl := &sk.Kennung

	if !kennung.Contains(idl, k) {
		return
	}

	ok = true

	return
}

func (k Kennung) String() string {
	var sb strings.Builder

	if k.Exact {
		sb.WriteRune('=')
	}

	if k.Virtual {
		sb.WriteRune('%')
	}

	sb.WriteString(kennung.FormattedString(k.Kennung2))

	return sb.String()
}

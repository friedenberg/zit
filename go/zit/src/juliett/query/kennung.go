package query

import (
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Kennung struct {
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
		var idx int

		if k.Exact {
			idx, ok = me.Verzeichnisse.Etiketten.All.ContainsKennungEtikettExact(
				k.Kennung2,
			)
		} else {
			idx, ok = me.Verzeichnisse.Etiketten.All.ContainsKennungEtikett(
				k.Kennung2,
			)
		}

		ui.Log().Print(k, idx, ok, me.Verzeichnisse.Etiketten.All, sk)

		if ok {
			// if k.Exact {
			// 	ewp := me.Verzeichnisse.Etiketten.All[idx]
			// 	ui.Debug().Print(ewp, sk)
			// }

			ps := me.Verzeichnisse.Etiketten.All[idx]
			sk.Metadatei.Verzeichnisse.QueryPath.Push(ps.Parents)
			return
		}

		return

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

package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type ObjekteOptions = sku.ObjekteOptions

func (s *Store) ReadOneKennungExternal(
	o ObjekteOptions,
	k1 schnittstellen.StringerGattungKastenGetter,
	sk *sku.Transacted,
) (el sku.ExternalLike, err error) {
	switch k1.GetKasten().GetKastenString() {
	case "chrome":
		// TODO populate with chrome kasten
		ui.Debug().Print("would populate from chrome")

	default:
		if el, err = s.cwdFiles.ReadKennung(o, k1, sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) Open(
	kasten schnittstellen.KastenGetter,
	m checkout_mode.Mode,
	ph schnittstellen.FuncIter[string],
	zsc sku.CheckedOutLikeSet,
) (err error) {
	switch kasten.GetKasten().GetKastenString() {
	case "chrome":
		err = todo.Implement()

	default:
		if err = s.cwdFiles.Open(m, ph, zsc); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

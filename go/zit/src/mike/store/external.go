package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
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
		if el, err = s.cwdFiles.ReadOneKennung(o, k1, sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) ReadOneCheckedOutFS(
	o ObjekteOptions,
	em *store_fs.KennungFDPair,
) (co *store_fs.CheckedOut, err error) {
	co = store_fs.GetCheckedOutPool().Get()

	if err = s.ReadOneInto(&em.Kennung, &co.Internal); err != nil {
		if collections.IsErrNotFound(err) {
			// TODO mark status as new
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.cwdFiles.ReadOneKennungFDPairExternalInto(
		o,
		em,
		&co.Internal,
		&co.External,
	); err != nil {
		if errors.Is(err, store_fs.ErrExternalHasConflictMarker) {
			err = nil
			co.State = checked_out_state.StateConflicted
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	co.DetermineState(false)

	return
}

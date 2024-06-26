package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type Status struct{}

func init() {
	registerCommandWithExternalQuery(
		"status",
		func(f *flag.FlagSet) CommandWithExternalQuery {
			c := &Status{}

			return c
		},
	)
}

func (c Status) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(gattung.TrueGattung()...)
}

func (c Status) ModifyBuilder(
	b *query.Builder,
) {
	b.WithHidden(nil)
}

func (c Status) RunWithExternalQuery(
	u *umwelt.Umwelt,
	eqwk sku.ExternalQueryWithKasten,
) (err error) {
	pcol := u.PrinterCheckedOutForKasten(eqwk.Kasten)

	if err = u.GetStore().QueryCheckedOut(
		eqwk,
		func(co sku.CheckedOutLike) (err error) {
			if err = pcol(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	qg := eqwk.Queryable.(*query.Group)

	if err = u.GetStore().QueryUnsure(
		qg,
		store.UnsureMatchOptions{
			UnsureMatchType: store.UnsureMatchTypeMetadateiSansTaiHistory | store.UnsureMatchTypeBezeichnung,
		},
		u.PrinterMatching(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = qg.GetExplicitCwdFDs().Each(
		u.GetStore().GetCwdFiles().MarkUnsureAkten,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := u.PrinterCheckedOutForKasten(eqwk.Kasten)

	if err = u.GetStore().QueryAllMatchingAkten(
		qg,
		u.GetStore().GetCwdFiles().GetUnsureAkten(),
		func(fd *fd.FD, z *sku.Transacted) (err error) {
			if z == nil {
				err = u.PrinterFileNotRecognized()(fd)
				return
			}

			os := sha.Make(z.GetObjekteSha())
			as := sha.Make(z.GetAkteSha())

			fr := store_fs.GetCheckedOutPool().Get()
			defer store_fs.GetCheckedOutPool().Put(fr)

			fr.State = checked_out_state.StateRecognized

			if err = fr.Internal.SetFromSkuLike(z); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = fr.External.SetFromSkuLike(z); err != nil {
				err = errors.Wrap(err)
				return
			}

			fr.External.FDs.Akte.ResetWith(fd)
			fr.External.SetAkteSha(as)

			if err = fr.External.SetObjekteSha(os); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = p(fr); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type Status struct{}

func init() {
	registerCommandWithQuery(
		"status",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Status{}

			return c
		},
	)
}

func (c Status) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(gattung.TrueGattung()...)
}

func (c Status) RunWithQuery(
	u *umwelt.Umwelt,
	qg *query.Group,
) (err error) {
	pcol := u.PrinterCheckedOut()

	if err = u.GetStore().ReadFiles(
		qg,
		func(co *sku.CheckedOut) (err error) {
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

	if err = qg.GetExplicitCwdFDs().Each(
		u.GetStore().GetCwdFiles().MarkUnsureAkten,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := u.PrinterCheckedOut()

	if err = u.GetStore().ReadAllMatchingAkten(
		qg,
		u.GetStore().GetCwdFiles().UnsureAkten,
		func(fd *fd.FD, z *sku.Transacted) (err error) {
			if z == nil {
				err = u.PrinterFileNotRecognized()(fd)
				return
			}

			os := sha.Make(z.GetObjekteSha())
			as := sha.Make(z.GetAkteSha())

			fr := sku.GetCheckedOutPool().Get()
			defer sku.GetCheckedOutPool().Put(fr)

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

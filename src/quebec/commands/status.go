package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/checked_out_state"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/checked_out"
	"github.com/friedenberg/zit/src/lima/cwd"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Status struct{}

func init() {
	registerCommandWithCwdQuery(
		"status",
		func(f *flag.FlagSet) CommandWithCwdQuery {
			c := &Status{}

			return c
		},
	)
}

func (c Status) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(gattung.TrueGattung()...)
}

func (c Status) RunWithCwdQuery(
	u *umwelt.Umwelt,
	ms matcher.Query,
	possible *cwd.CwdFiles,
) (err error) {
	pcol := u.PrinterCheckedOutLike()

	if err = u.StoreObjekten().ReadFiles(
		possible,
		objekte.MakeFuncReaderTransactedLikePtr(ms, u.StoreObjekten().Query),
		iter.MakeChain(
			objekte.MakeFilterFromMetaSet(ms),
			func(co *sku.CheckedOut) (err error) {
				if err = pcol(co); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := u.PrinterCheckedOutLike()

	if err = u.StoreObjekten().ReadAllMatchingAkten(
		possible.UnsureAkten,
		func(fd kennung.FD, z sku.SkuLikePtr) (err error) {
			if z == nil {
				err = u.PrinterFileNotRecognized()(&fd)
			} else {
				os := sha.Make(z.GetObjekteSha())
				as := sha.Make(z.GetAkteSha())

				fr := &checked_out.Zettel{
					State: checked_out_state.StateRecognized,
				}

				if err = fr.Internal.SetFromSkuLike(z); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = fr.External.SetFromSkuLike(z); err != nil {
					err = errors.Wrap(err)
					return
				}

				fr.External.FDs = sku.ExternalFDs{
					Akte: fd,
				}
				fr.External.SetAkteSha(as)
				fr.External.ObjekteSha = os

				err = p(fr)
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

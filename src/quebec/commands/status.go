package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/november/umwelt"
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
	ms kennung.MetaSet,
	possible cwd.CwdFiles,
) (err error) {
	pcol := u.PrinterCheckedOutLike()

	if err = u.StoreWorkingDirectory().ReadFiles(
		possible,
		ms,
		iter.MakeChain(
			objekte.MakeFilterFromMetaSet(ms),
			func(co objekte.CheckedOutLike) (err error) {
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
		func(fd kennung.FD, z *zettel.Transacted) (err error) {
			if z == nil {
				err = u.PrinterFileNotRecognized()(&fd)
			} else {
				os := sha.Make(z.GetObjekteSha())
				as := sha.Make(z.GetAkteSha())

				fr := &zettel.CheckedOut{
					State:    objekte.CheckedOutStateRecognized,
					Internal: *z,
					External: zettel.External{
						Objekte: z.Akte,
						Sku: sku.External[kennung.Hinweis, *kennung.Hinweis]{
							Kennung: z.Sku.Kennung,
							FDs: sku.ExternalFDs{
								Akte: fd,
							},
						},
					},
				}

				fr.External.SetAkteSha(as)
				fr.External.Sku.ObjekteSha = os

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

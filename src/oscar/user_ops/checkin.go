package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Checkin struct {
	Delete bool
}

func (c Checkin) Run(
	u *umwelt.Umwelt,
	ms kennung.MetaSet,
	pz cwd.CwdFiles,
) (err error) {
	fds := kennung.MakeMutableFDSet()

	u.Lock()
	defer errors.Deferred(&err, u.Unlock)

	if err = u.StoreWorkingDirectory().ReadFiles(
		pz,
		ms,
		iter.MakeChain(
			objekte.MakeFilterFromMetaSet(ms),
			func(co objekte.CheckedOutLike) (err error) {
				// var tl objekte.TransactedLike

				switch aco := co.(type) {
				case *zettel.CheckedOut:
					if _, err = u.StoreObjekten().Zettel().UpdateCheckedOut(
						*aco,
					); err != nil {
						err = errors.Wrap(err)
						return
					}

				case *typ.CheckedOut:
					if _, err = u.StoreObjekten().Typ().CreateOrUpdateCheckedOut(
						aco,
					); err != nil {
						err = errors.Wrap(err)
						return
					}

				case *etikett.CheckedOut:
					if _, err = u.StoreObjekten().Etikett().CreateOrUpdateCheckedOut(
						aco,
					); err != nil {
						err = errors.Wrap(err)
						return
					}

				default:
					err = errors.Implement()
					return
				}

				e := co.GetExternal()
				fds.Add(e.GetObjekteFD())
				fds.Add(e.GetAkteFD())

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !c.Delete {
		return
	}

	deleteOp := DeleteCheckout{
		Umwelt: u,
	}

	if err = deleteOp.Run(fds); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

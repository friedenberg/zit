package user_ops

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/checked_out"
	"github.com/friedenberg/zit/src/lima/cwd"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Checkin struct {
	Delete bool
}

func (c Checkin) Run(
	u *umwelt.Umwelt,
	ms matcher.Query,
	pz *cwd.CwdFiles,
) (err error) {
	fds := collections_ptr.MakeMutableValueSet[kennung.FD, *kennung.FD](nil)
	l := &sync.Mutex{}

	u.Lock()
	defer errors.Deferred(&err, u.Unlock)

	if err = u.StoreWorkingDirectory().ReadFiles(
		pz,
		ms,
		iter.MakeChain(
			objekte.MakeFilterFromMetaSet(ms),
			func(co objekte.CheckedOutLikePtr) (err error) {
				switch aco := co.(type) {
				case *checked_out.Zettel:
					if _, err = u.StoreObjekten().Zettel().UpdateCheckedOut(
						aco,
					); err != nil {
						err = errors.Wrap(err)
						return
					}

				case *checked_out.Kasten:
					if _, err = u.StoreObjekten().Kasten().CreateOrUpdateCheckedOut(
						aco,
					); err != nil {
						err = errors.Wrap(err)
						return
					}

				case *objekte.CheckedOut2:
					if _, err = u.StoreObjekten().CreateOrUpdator.CreateOrUpdateCheckedOut(
						aco,
					); err != nil {
						err = errors.Wrap(err)
						return
					}

				case *checked_out.Typ:
					if _, err = u.StoreObjekten().Typ().CreateOrUpdateCheckedOut(
						aco,
					); err != nil {
						err = errors.Wrap(err)
						return
					}

				case *checked_out.Etikett:
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

				e := co.GetExternalLike()

				l.Lock()
				defer l.Unlock()
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

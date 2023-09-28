package user_ops

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/oscar/umwelt"
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

	if err = u.StoreObjekten().ReadFiles(
		pz,
		matcher.MakeFuncReaderTransactedLikePtr(ms, u.StoreObjekten().Query),
		iter.MakeChain(
			matcher.MakeFilterFromQuery(ms),
			func(co *sku.CheckedOut) (err error) {
				if _, err = u.StoreObjekten().CreateOrUpdateCheckedOut(
					co,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				l.Lock()
				defer l.Unlock()

				fds.Add(co.External.GetObjekteFD())
				fds.Add(co.External.GetAkteFD())

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

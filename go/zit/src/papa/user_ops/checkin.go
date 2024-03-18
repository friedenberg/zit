package user_ops

import (
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/hotel/matcher_proto"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type Checkin struct {
	Delete bool
}

func (c Checkin) Run(
	u *umwelt.Umwelt,
	ms matcher_proto.QueryGroup,
) (err error) {
	fds := fd.MakeMutableSet()
	l := &sync.Mutex{}

	u.Lock()
	defer errors.Deferred(&err, u.Unlock)

	log.Log().Print(ms)

	if err = u.StoreObjekten().ReadFiles(
		matcher_proto.MakeFuncReaderTransactedLikePtr(ms, u.StoreObjekten().QueryWithoutCwd),
		iter.MakeChain(
			matcher_proto.MakeFilterFromQuery(ms),
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

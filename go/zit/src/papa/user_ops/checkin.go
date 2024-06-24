package user_ops

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type Checkin struct {
	Delete bool
}

func (c Checkin) Run(
	u *umwelt.Umwelt,
	qg query.GroupWithKasten,
) (err error) {
	fds := fd.MakeMutableSet()
	l := &sync.Mutex{}

	u.Lock()
	defer errors.Deferred(&err, u.Unlock)

	ui.Log().Print(qg)

	if err = u.GetStore().ReadExternal(
		qg,
		func(col sku.CheckedOutLike) (err error) {
			ui.Log().Print(col)

			switch cot := col.(type) {
			default:
				err = todo.Implement()

			case *store_fs.CheckedOut:
				if _, err = u.GetStore().CreateOrUpdateCheckedOut(
					cot,
					true,
				); err != nil {
					ui.Debug().Print(err)
					err = errors.Wrap(err)
					return
				}

				l.Lock()
				defer l.Unlock()

				// TODO support generic deletes
				fds.Add(cot.External.GetObjekteFD())
				fds.Add(cot.External.GetAkteFD())
			}

			return
		},
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

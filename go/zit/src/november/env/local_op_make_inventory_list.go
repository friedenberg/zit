package env

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func (u *Local) MakeInventoryList(
	qg *query.Group,
) (list *sku.List, err error) {
	list = sku.MakeList()

	var l sync.Mutex

	if err = u.GetStore().QueryTransacted(
		qg,
		func(sk *sku.Transacted) (err error) {
			l.Lock()
			defer l.Unlock()

			list.Add(sk.CloneTransacted())

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

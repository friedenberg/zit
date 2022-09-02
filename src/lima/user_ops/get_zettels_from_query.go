package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/hotel/collections"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
	"github.com/friedenberg/zit/src/hotel/zettel_named"
	"github.com/friedenberg/zit/zettel_transacted"
)

type GetZettelsFromQuery struct {
	Umwelt *umwelt.Umwelt
}

func (c GetZettelsFromQuery) Run(query zettel_named.NamedFilter) (result collections.SetTransacted, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	var set map[hinweis.Hinweis]zettel_transacted.Transacted

	if set, err = store.StoreObjekten().ZettelenSchwanzen(query); err != nil {
		err = errors.Error(err)
		return
	}

	result = collections.MakeSetUniqueTransacted(len(set))

	for _, tz := range set {
		result.Add(tz)
	}

	return
}

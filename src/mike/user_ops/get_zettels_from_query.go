package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/hotel/zettel_named"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/lima/store_with_lock"
)

type GetZettelsFromQuery struct {
	Umwelt *umwelt.Umwelt
}

func (c GetZettelsFromQuery) Run(query zettel_named.NamedFilter) (result zettel_transacted.Set, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	var set map[hinweis.Hinweis]zettel_transacted.Zettel

	if set, err = store.StoreObjekten().ZettelenSchwanzen(query); err != nil {
		err = errors.Wrap(err)
		return
	}

	result = zettel_transacted.MakeSetUnique(len(set))

	for _, tz := range set {
		result.Add(tz)
	}

	return
}

package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	zettel_stored "github.com/friedenberg/zit/src/golf/zettel_stored"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
)

type GetZettelsFromQuery struct {
	Umwelt *umwelt.Umwelt
}

func (c GetZettelsFromQuery) Run(query zettel_stored.NamedFilter) (result ZettelResults, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	var set map[hinweis.Hinweis]zettel_stored.Transacted

	if set, err = store.StoreObjekten().ZettelenSchwanzen(query); err != nil {
		err = errors.Error(err)
		return
	}

	result.SetNamed = zettel_stored.MakeSetNamed()

	for h, tz := range set {
		result.SetNamed[h.String()] = tz.Named
	}

	return
}

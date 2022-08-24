package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/india/store_with_lock"
)

type GetZettelsFromQuery struct {
	Umwelt *umwelt.Umwelt
}

func (c GetZettelsFromQuery) Run(query stored_zettel.NamedFilter) (result ZettelResults, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	var set map[hinweis.Hinweis]stored_zettel.Transacted

	if set, err = store.Zettels().ZettelenSchwanzen(query); err != nil {
		err = errors.Error(err)
		return
	}

	result.SetNamed = stored_zettel.MakeSetNamed()

	for h, tz := range set {
		result.SetNamed[h.String()] = tz.Named
	}

	return
}

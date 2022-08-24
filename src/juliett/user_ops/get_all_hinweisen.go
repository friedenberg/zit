package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/india/store_with_lock"
)

type GetAllHinweisen struct {
	Umwelt *umwelt.Umwelt
}

type GetAllHinweisenResults struct {
	Hinweisen        []hinweis.Hinweis
	HinweisenStrings []string
}

func (op GetAllHinweisen) Run() (results GetAllHinweisenResults, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(op.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	var zs map[hinweis.Hinweis]stored_zettel.Transacted

	if zs, err = store.Zettels().ZettelenSchwanzen(); err != nil {
		err = errors.Error(err)
		return
	}

	results.Hinweisen = make([]hinweis.Hinweis, len(zs))
	results.HinweisenStrings = make([]string, len(zs))

	i := 0

	for h, _ := range zs {
		results.Hinweisen[i] = h
		results.HinweisenStrings[i] = h.String()
		i++
	}

	return
}

package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/store_with_lock"
	"github.com/friedenberg/zit/src/juliett/umwelt"
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
		err = errors.Wrap(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	var zs map[hinweis.Hinweis]zettel_transacted.Zettel

	if zs, err = store.StoreObjekten().ZettelenSchwanzen(); err != nil {
		err = errors.Wrap(err)
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

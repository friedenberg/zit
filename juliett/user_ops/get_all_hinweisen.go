package user_ops

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/india/store_with_lock"
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

	if _, results.Hinweisen, err = store.Hinweisen().All(); err != nil {
		err = _Error(err)
		return
	}

	results.HinweisenStrings = make([]string, len(results.Hinweisen))

	for i, h := range results.Hinweisen {
		results.HinweisenStrings[i] = h.String()
	}

	return
}

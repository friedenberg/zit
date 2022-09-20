package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type GetAllHinweisen struct {
	*umwelt.Umwelt
}

type GetAllHinweisenResults struct {
	Hinweisen        []hinweis.Hinweis
	HinweisenStrings []string
}

func (op GetAllHinweisen) Run() (results GetAllHinweisenResults, err error) {
	var zts zettel_transacted.Set

	if zts, err = op.StoreObjekten().ZettelenSchwanzen(); err != nil {
		err = errors.Wrap(err)
		return
	}

	results.Hinweisen = make([]hinweis.Hinweis, zts.Len())
	results.HinweisenStrings = make([]string, zts.Len())

	i := 0

	zts.Each(
		func(zt zettel_transacted.Zettel) (err error) {
			h := zt.Named.Hinweis
			results.Hinweisen[i] = h
			results.HinweisenStrings[i] = h.String()
			i++

			return
		},
	)

	return
}

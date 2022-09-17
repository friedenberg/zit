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
	var zs map[hinweis.Hinweis]zettel_transacted.Zettel

	if zs, err = op.StoreObjekten().ZettelenSchwanzen(); err != nil {
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

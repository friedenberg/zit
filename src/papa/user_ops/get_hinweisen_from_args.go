package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/golf/hinweis"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type GetHinweisenFromArgs struct {
	*umwelt.Umwelt
}

func (u GetHinweisenFromArgs) RunOne(v string) (h hinweis.Hinweis, err error) {
	if h, err = u.StoreObjekten().Abbr().ExpandHinweisString(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (u GetHinweisenFromArgs) RunMany(vs ...string) (hs []hinweis.Hinweis, err error) {
	hs = make([]hinweis.Hinweis, len(vs))

	for i, _ := range hs {
		if hs[i], err = u.RunOne(vs[i]); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

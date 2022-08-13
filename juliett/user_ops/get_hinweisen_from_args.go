package user_ops

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/charlie/hinweis"
)

type GetHinweisenFromArgs struct {
}

func (u GetHinweisenFromArgs) RunOne(v string) (h hinweis.Hinweis, err error) {
	if h, err = hinweis.MakeBlindHinweis(v); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (u GetHinweisenFromArgs) RunMany(vs ...string) (hs []hinweis.Hinweis, err error) {
	hs = make([]hinweis.Hinweis, len(vs))

	for i, _ := range hs {
		if hs[i], err = u.RunOne(vs[i]); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}

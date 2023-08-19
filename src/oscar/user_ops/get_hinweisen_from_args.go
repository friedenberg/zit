package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type GetHinweisenFromArgs struct {
	*umwelt.Umwelt
}

func (u GetHinweisenFromArgs) RunOne(v string) (h kennung.Hinweis, err error) {
	if h, err = u.StoreObjekten().GetAbbrStore().Hinweis().ExpandString(
		v,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (u GetHinweisenFromArgs) RunMany(vs ...string) (hs []kennung.Hinweis, err error) {
	hs = make([]kennung.Hinweis, len(vs))

	for i := range hs {
		if hs[i], err = u.RunOne(vs[i]); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

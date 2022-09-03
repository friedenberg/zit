package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/delta/umwelt"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/juliett/store_with_lock"
)

type Checkout struct {
	*umwelt.Umwelt
	store_working_directory.CheckoutOptions
}

func (c Checkout) RunManyHinweisen(
	s store_with_lock.Store,
	hins ...hinweis.Hinweis,
) (results zettel_checked_out.Set, err error) {
	zts := zettel_transacted.MakeSetUnique(len(hins))

	for _, h := range hins {
		var zt zettel_transacted.Zettel

		if zt, err = s.StoreObjekten().Read(h); err != nil {
			err = errors.Wrap(err)
			return
		}

		zts.Add(zt)
	}

	results = zettel_checked_out.MakeSetUnique(zts.Len())
	ztsl := zts.ToSlice()

	if results, err = s.StoreWorkingDirectory().Checkout(c.CheckoutOptions, ztsl...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

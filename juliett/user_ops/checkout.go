package user_ops

import (
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/delta/hinweis"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	checkout_store "github.com/friedenberg/zit/golf/store_checkout"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type Checkout struct {
	Options checkout_store.CheckinOptions
	Umwelt  *umwelt.Umwelt
}

type CheckoutResults struct {
	Zettelen      []stored_zettel.CheckedOut
	FilesZettelen []string
	FilesAkten    []string
}

func (c Checkout) RunManyHinweisen(
	s store_with_lock.Store,
	hins ...hinweis.Hinweis,
) (results CheckoutResults, err error) {
	zs := make([]stored_zettel.Transacted, len(hins))

	for i, _ := range zs {
		h := hins[i]

		if zs[i], err = s.Zettels().Read(h); err != nil {
			err = errors.Error(err)
			return
		}
	}

	if results.Zettelen, err = s.CheckoutStore().Checkout(c.Options, zs...); err != nil {
		err = errors.Error(err)
		return
	}

	results.FilesZettelen = make([]string, 0, len(results.Zettelen))
	results.FilesAkten = make([]string, 0)

	for _, z := range results.Zettelen {
		results.FilesZettelen = append(results.FilesZettelen, z.External.Path)

		if z.External.AktePath != "" {
			results.FilesAkten = append(results.FilesAkten, z.External.AktePath)
		}
	}

	return
}

package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	store_checkout "github.com/friedenberg/zit/src/hotel/store_checkout"
	"github.com/friedenberg/zit/src/india/zettel_checked_out"
	"github.com/friedenberg/zit/src/india/store_with_lock"
)

type Checkout struct {
	*umwelt.Umwelt
	store_checkout.CheckoutOptions
}

type CheckoutResults struct {
	Zettelen      []zettel_checked_out.CheckedOut
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

	if results.Zettelen, err = s.CheckoutStore().Checkout(c.CheckoutOptions, zs...); err != nil {
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

package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	zettel_stored "github.com/friedenberg/zit/src/golf/zettel_stored"
	"github.com/friedenberg/zit/src/india/zettel_checked_out"
	store_working_directory "github.com/friedenberg/zit/src/juliett/store_working_directory"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
)

type Checkout struct {
	*umwelt.Umwelt
	store_working_directory.CheckoutOptions
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
	zs := make([]zettel_stored.Transacted, len(hins))

	for i, _ := range zs {
		h := hins[i]

		if zs[i], err = s.StoreObjekten().Read(h); err != nil {
			err = errors.Error(err)
			return
		}
	}

	if results.Zettelen, err = s.StoreWorkingDirectory().Checkout(c.CheckoutOptions, zs...); err != nil {
		err = errors.Error(err)
		return
	}

	results.FilesZettelen = make([]string, 0, len(results.Zettelen))
	results.FilesAkten = make([]string, 0)

	for _, z := range results.Zettelen {
		results.FilesZettelen = append(results.FilesZettelen, z.External.ZettelFD.Path)

		if z.External.AkteFD.Path != "" {
			results.FilesAkten = append(results.FilesAkten, z.External.AkteFD.Path)
		}
	}

	return
}

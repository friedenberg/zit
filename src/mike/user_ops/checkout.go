package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/kilo/store_working_directory"
	"github.com/friedenberg/zit/src/lima/store_with_lock"
)

type Checkout struct {
	*umwelt.Umwelt
	store_working_directory.CheckoutOptions
}

type CheckoutResults struct {
	Zettelen      []zettel_checked_out.Zettel
	FilesZettelen []string
	FilesAkten    []string
}

func (c Checkout) RunManyHinweisen(
	s store_with_lock.Store,
	hins ...hinweis.Hinweis,
) (results CheckoutResults, err error) {
	zs := make([]zettel_transacted.Zettel, len(hins))

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

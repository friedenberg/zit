package user_ops

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type Checkout struct {
	Options _ZettelsCheckinOptions
	Umwelt  *umwelt.Umwelt
}

type CheckoutResults struct {
	Zettelen      []_ZettelCheckedOut
	FilesZettelen []string
	FilesAkten    []string
}

func (c Checkout) Run(args ...string) (results CheckoutResults, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	if results.Zettelen, err = store.Zettels().Checkout(c.Options, args...); err != nil {
		err = _Error(err)
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

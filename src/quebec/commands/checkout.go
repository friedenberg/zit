package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/lima/cwd"
	"github.com/friedenberg/zit/src/mike/store_util"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Checkout struct {
	store_util.CheckoutOptions
}

func init() {
	registerCommandWithQuery(
		"checkout",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Checkout{}

			c.CheckoutOptions.AddToFlagSet(f)

			return c
		},
	)
}

func (c Checkout) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
	)
}

func (c Checkout) RunWithQuery(
	u *umwelt.Umwelt,
	ms matcher.Query,
) (err error) {
	var cwdFiles cwd.CwdFiles

	if cwdFiles, err = cwd.MakeCwdFilesAll(
		u.Konfig(),
		u.Standort().Cwd(),
		u.Standort(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.CheckoutOptions.Cwd = cwdFiles

	if err = u.StoreObjekten().CheckoutQuery(
		c.CheckoutOptions,
		matcher.MakeFuncReaderTransactedLikePtr(ms, u.StoreObjekten().Query),
		func(co *sku.CheckedOut) (err error) {
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

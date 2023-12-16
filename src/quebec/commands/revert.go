package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Revert struct{}

func init() {
	registerCommandWithQuery(
		"revert",
		func(_ *flag.FlagSet) CommandWithQuery {
			c := &Revert{}

			return c
		},
	)
}

func (c Revert) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
		gattung.Bestandsaufnahme,
		gattung.Kasten,
	)
}

func (c Revert) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
		// gattung.Bestandsaufnahme,
		gattung.Kasten,
	)
}

func (c Revert) RunWithQuery(u *umwelt.Umwelt, ms matcher.Query) (err error) {
	kinderToMutter := make(map[string]string)

	if err = u.StoreObjekten().QueryWithCwd(
		ms,
		func(z *sku.Transacted) (err error) {
			mu := z.Metadatei.Verzeichnisse.Mutter

			if mu.IsNull() {
				log.Err().Printf("%s has null mutter, cannot revert", z)
				return
			}

			sh := z.Metadatei.Verzeichnisse.Sha

			kinderToMutter[sh.String()] = mu.String()

			log.Debug().Printf("%s -> %s", sh, mu)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

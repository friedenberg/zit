package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	checkout_store "github.com/friedenberg/zit/src/hotel/store_checkout"
	"github.com/friedenberg/zit/src/india/store_with_lock"
)

type Checkin struct {
	*umwelt.Umwelt
	checkout_store.OptionsReadExternal
}

type CheckinResults struct {
	Zettelen map[hinweis.Hinweis]stored_zettel.CheckedOut
}

func (c Checkin) Run(
	store store_with_lock.Store,
	zettelen ...stored_zettel.External,
) (results CheckinResults, err error) {
	results.Zettelen = make(map[hinweis.Hinweis]stored_zettel.CheckedOut)

	for _, z := range zettelen {
		var tz stored_zettel.Transacted

		if tz, err = store.Zettels().Update(z.Hinweis, z.Zettel); err != nil {
			err = errors.Error(err)
			return
		}

		//TODO: add states to checkin process to indicate results of update call
		// stdprinter.Outf("%s (unchanged)\n", tz.Named)
		logz.PrintDebug(tz)
		stdprinter.Outf("%s (updated)\n", tz.Named)

		results.Zettelen[tz.Hinweis] = stored_zettel.CheckedOut{
			Internal: tz,
			External: z,
		}
	}

	return
}

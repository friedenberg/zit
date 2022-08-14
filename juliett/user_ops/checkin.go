package user_ops

import (
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/delta/hinweis"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	checkout_store "github.com/friedenberg/zit/golf/store_checkout"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type Checkin struct {
	Umwelt  *umwelt.Umwelt
	Options checkout_store.CheckinOptions
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
		stdprinter.Outf("%s (updated)\n", tz.Named)

		results.Zettelen[tz.Hinweis] = stored_zettel.CheckedOut{
			Internal: tz,
			External: z,
		}
	}

	return
}

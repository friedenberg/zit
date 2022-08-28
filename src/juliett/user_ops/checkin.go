package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	zettel_stored "github.com/friedenberg/zit/src/golf/zettel_stored"
	store_working_directory "github.com/friedenberg/zit/src/hotel/store_working_directory"
	"github.com/friedenberg/zit/src/india/store_with_lock"
	"github.com/friedenberg/zit/src/india/zettel_checked_out"
)

type Checkin struct {
	*umwelt.Umwelt
	store_working_directory.OptionsReadExternal
}

type CheckinResults struct {
	Zettelen map[hinweis.Hinweis]zettel_checked_out.CheckedOut
}

func (c Checkin) Run(
	store store_with_lock.Store,
	zettelen ...zettel_stored.External,
) (results CheckinResults, err error) {
	results.Zettelen = make(map[hinweis.Hinweis]zettel_checked_out.CheckedOut)

	for _, z := range zettelen {
		var tz zettel_stored.Transacted

		if tz, err = store.StoreObjekten().Update(z.Hinweis, z.Named.Stored.Zettel); err != nil {
			err = errors.Error(err)
			return
		}

		//TODO: add states to checkin process to indicate results of update call
		// stdprinter.Outf("%s (unchanged)\n", tz.Named)
		logz.PrintDebug(tz)
		stdprinter.Outf("%s (updated)\n", tz.Named)

		results.Zettelen[tz.Named.Hinweis] = zettel_checked_out.CheckedOut{
			Internal: tz,
			External: z,
		}
	}

	return
}

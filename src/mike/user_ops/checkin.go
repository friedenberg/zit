package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/kilo/store_working_directory"
	"github.com/friedenberg/zit/src/lima/store_with_lock"
)

type Checkin struct {
	*umwelt.Umwelt
	store_working_directory.OptionsReadExternal
}

type CheckinResults struct {
	Zettelen map[hinweis.Hinweis]zettel_checked_out.Zettel
}

func (c Checkin) Run(
	store store_with_lock.Store,
	zettelen ...zettel_external.Zettel,
) (results CheckinResults, err error) {
	results.Zettelen = make(map[hinweis.Hinweis]zettel_checked_out.Zettel)

	for _, z := range zettelen {
		var tz zettel_transacted.Zettel

		if tz, err = store.StoreObjekten().Update(z.Named.Hinweis, z.Named.Stored.Zettel); err != nil {
			err = errors.Error(err)
			return
		}

		//TODO: add states to checkin process to indicate results of update call
		// stdprinter.Outf("%s (unchanged)\n", tz.Named)
		errors.PrintDebug(tz)
		stdprinter.Outf("%s (updated)\n", tz.Named)

		results.Zettelen[tz.Named.Hinweis] = zettel_checked_out.Zettel{
			Internal: tz,
			External: z,
		}
	}

	return
}

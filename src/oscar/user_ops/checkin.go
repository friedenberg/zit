package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	store_fs "github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Checkin struct {
	*umwelt.Umwelt
	store_fs.OptionsReadExternal
}

type CheckinResults struct {
	Zettelen map[hinweis.Hinweis]zettel_checked_out.Zettel
}

func (c Checkin) Run(
	zettelen ...zettel_external.Zettel,
) (results CheckinResults, err error) {
	results.Zettelen = make(map[hinweis.Hinweis]zettel_checked_out.Zettel)

	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer c.Unlock()

	for _, z := range zettelen {
		var tz zettel_transacted.Zettel

		if tz, err = c.StoreObjekten().Update(&z.Named); err != nil {
			err = errors.Wrap(err)
			return
		}

		//TODO: add states to checkin process to indicate results of update call
		// stdprinter.Outf("%s (unchanged)", tz.Named)

		results.Zettelen[tz.Named.Hinweis] = zettel_checked_out.Zettel{
			Internal: tz,
			External: z,
		}
	}

	return
}

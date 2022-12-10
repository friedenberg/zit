package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/mike/zettel_checked_out"
	"github.com/friedenberg/zit/src/november/umwelt"
)

// TODO-P4 move to store_fs
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
		var tz zettel.Transacted

		if tz, err = c.StoreObjekten().Zettel().Update(
			&z.Objekte,
			&z.Sku.Kennung,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		//TODO: add states to checkin process to indicate results of update call
		// stdprinter.Outf("%s (unchanged)", tz.Named)

		results.Zettelen[tz.Sku.Kennung] = zettel_checked_out.Zettel{
			Internal: tz,
			External: z,
		}
	}

	return
}

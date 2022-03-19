package user_ops

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type Checkin struct {
	Umwelt  *umwelt.Umwelt
	Options _ZettelsCheckinOptions
}

type CheckinResults struct {
	Zettelen map[hinweis.Hinweis]stored_zettel.CheckedOut
}

func (c Checkin) Run(zettelen ...stored_zettel.External) (results CheckinResults, err error) {
	results.Zettelen = make(map[hinweis.Hinweis]stored_zettel.CheckedOut)

	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	for _, z := range zettelen {
		named := stored_zettel.Named{
			Hinweis: z.Hinweis,
			Stored: stored_zettel.Stored{
				Zettel: z.Zettel,
			},
		}

		if named, err = store.Zettels().Update(named); err != nil {
			err = _Error(err)
			return
		}

		results.Zettelen[named.Hinweis] = stored_zettel.CheckedOut{
			Internal: named,
			External: z,
		}
	}

	return
}

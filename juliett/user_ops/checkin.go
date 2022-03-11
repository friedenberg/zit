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

func (c Checkin) Run(args ...string) (results CheckinResults, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	if results.Zettelen, err = store.Zettels().Checkin(c.Options, args...); err != nil {
		err = _Error(err)
		return
	}

	return
}

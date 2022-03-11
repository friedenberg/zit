package user_ops

import (
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type Checkin struct {
	Umwelt  _Umwelt
	Store   _Store
	Options _ZettelsCheckinOptions
}

type CheckinResults struct {
	Zettelen map[hinweis.Hinweis]stored_zettel.CheckedOut
}

func (c Checkin) Run(args ...string) (results CheckinResults, err error) {
	if results.Zettelen, err = c.Store.Checkin(c.Options, args...); err != nil {
		err = _Error(err)
		return
	}

	return
}

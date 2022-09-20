package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type GetZettelsFromQuery struct {
	*umwelt.Umwelt
}

func (c GetZettelsFromQuery) Run(query zettel_named.NamedFilter) (result zettel_transacted.Set, err error) {
	if result, err = c.StoreObjekten().ZettelenSchwanzen(query); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

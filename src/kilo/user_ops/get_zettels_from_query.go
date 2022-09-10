package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/umwelt"
)

type GetZettelsFromQuery struct {
	*umwelt.Umwelt
}

func (c GetZettelsFromQuery) Run(query zettel_named.NamedFilter) (result zettel_transacted.Set, err error) {
	var set map[hinweis.Hinweis]zettel_transacted.Zettel

	if set, err = c.StoreObjekten().ZettelenSchwanzen(query); err != nil {
		err = errors.Wrap(err)
		return
	}

	result = zettel_transacted.MakeSetUnique(len(set))

	for _, tz := range set {
		result.Add(tz)
	}

	return
}

package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Pull struct{}

func init() {
	registerCommandWithRemoteAndQuery(
		"pull",
		func(f *flag.FlagSet) CommandWithRemoteAndQuery {
			c := &Pull{}

			return c
		},
	)
}

func (c Pull) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c Pull) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
	// return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Pull) RunWithRemoteAndQuery(
	local *env.Local,
	remote env.Env,
	qg *query.Group,
) {
	if err := local.PullQueryGroupFromRemote(
		remote,
		qg,
		true,
	); err != nil {
		local.Context.Cancel(errors.Wrap(err))
		return
	}
}

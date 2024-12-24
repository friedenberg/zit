package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Clone struct {
	env.BigBang
}

func init() {
	registerCommandWithRemoteAndQueryAndWithoutEnvironment(
		"clone",
		func(f *flag.FlagSet) CommandWithRemoteAndQuery {
			c := &Clone{
				BigBang: env.BigBang{
					Config:             immutable_config.Default(),
					ExcludeDefaultType: true,
				},
			}

			c.BigBang.AddToFlagSet(f)

			return c
		},
	)
}

func (c Clone) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c Clone) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
	// return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Clone) RunWithRemoteAndQuery(
	local *env.Local,
	remote env.Env,
	qg *query.Group,
) {
	if err := local.Start(c.BigBang); err != nil {
		local.CancelWithError(err)
		return
	}

	if err := local.PullQueryGroupFromRemote(
		remote,
		qg,
		true,
	); err != nil {
		local.CancelWithError(err)
		return
	}

	return
}

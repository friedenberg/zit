package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Test struct{}

func init() {
	registerCommandWithQuery(
		"test",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Test{}

			return c
		},
	)
}

func (c Test) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
		genres.Tag,
		genres.Type,
		genres.InventoryList,
		genres.Repo,
	)
}

func (c Test) DefaultGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
		genres.Tag,
		genres.Type,
		// gattung.Bestandsaufnahme,
		genres.Repo,
	)
}

func (c Test) RunWithQuery(
	u *env.Env,
	ms *query.Group,
) (err error) {
	u.GetStore().GetStreamIndex().GetProbeIndex().PrintAll()
	return
}

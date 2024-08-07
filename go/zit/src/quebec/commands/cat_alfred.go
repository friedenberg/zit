package commands

import (
	"bufio"
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/india/alfred"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type CatAlfred struct {
	Command
}

func init() {
	registerCommandWithQuery(
		"cat-alfred",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &CatAlfred{}

			return c
		},
	)
}

func (c CatAlfred) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Tag,
		genres.Zettel,
		genres.Type,
	)
}

func (c CatAlfred) DefaultGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Tag,
		genres.Zettel,
		genres.Type,
	)
}

func (c CatAlfred) RunWithQuery(
	u *env.Env,
	qg *query.Group,
) (err error) {
	// this command does its own error handling
	wo := bufio.NewWriter(u.Out())
	defer errors.DeferredFlusher(&err, wo)

	var aw *alfred.Writer

	if aw, err = alfred.New(
		wo,
		u.GetStore().GetAbbrStore().GetAbbr(),
		u.SkuFmtOrganize(qg.RepoId),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if err = u.GetStore().QueryWithKasten(
		qg,
		aw.PrintOne,
	); err != nil {
		aw.WriteError(err)
		return
	}

	return
}

package commands

import (
	"bufio"
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/india/alfred"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/oscar/umwelt"
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

func (c CatAlfred) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Etikett,
		gattung.Zettel,
		gattung.Typ,
	)
}

func (c CatAlfred) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Etikett,
		gattung.Zettel,
		gattung.Typ,
	)
}

func (c CatAlfred) RunWithQuery(
	u *umwelt.Umwelt,
	ms matcher.Query,
) (err error) {
	// this command does its own error handling
	wo := bufio.NewWriter(u.Out())
	defer errors.DeferredFlusher(&err, wo)

	var aw *alfred.Writer

	if aw, err = alfred.New(
		wo,
		u.StoreObjekten().GetKennungIndex(),
		u.StoreObjekten().GetKennungIndex(),
		u.StoreObjekten().GetAbbrStore().Hinweis().Abbreviate,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if err = u.StoreObjekten().QueryWithCwd(
		ms,
		aw.PrintOne,
	); err != nil {
		aw.WriteError(err)
		return
	}

	return
}

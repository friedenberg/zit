package commands

import (
	"bufio"
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/india/alfred"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
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

func (c CatAlfred) CompletionGattung() kennung.Gattung {
	return kennung.MakeGattung(
		gattung.Etikett,
		gattung.Zettel,
		gattung.Typ,
	)
}

func (c CatAlfred) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(
		gattung.Etikett,
		gattung.Zettel,
		gattung.Typ,
	)
}

func (c CatAlfred) RunWithQuery(
	u *umwelt.Umwelt,
	ms *query.Group,
) (err error) {
	// this command does its own error handling
	wo := bufio.NewWriter(u.Out())
	defer errors.DeferredFlusher(&err, wo)

	var aw *alfred.Writer

	if aw, err = alfred.New(
		wo,
		u.GetStore().GetAbbrStore().GetAbbr(),
		u.SkuFmtOrganize(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if err = u.GetStore().QueryWithKasten(
		query.GroupWithKasten{Group: ms},
		aw.PrintOne,
	); err != nil {
		aw.WriteError(err)
		return
	}

	return
}

package commands

import (
	"bufio"
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
	"github.com/friedenberg/zit/src/kilo/alfred"
	"github.com/friedenberg/zit/src/november/umwelt"
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

func (c CatAlfred) RunWithQuery(u *umwelt.Umwelt, ms kennung.MetaSet) (err error) {
	// this command does its own error handling
	wo := bufio.NewWriter(u.Out())
	defer errors.DeferredFlusher(&err, wo)

	var aw *alfred.Writer

	var ti kennung_index.KennungIndex[kennung.Typ]

	if ti, err = u.StoreObjekten().GetTypenIndex(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if aw, err = alfred.New(
		wo,
		u.StoreObjekten().GetKennungIndex(),
		ti,
		u.StoreObjekten().GetAbbrStore().Hinweis().Abbreviate,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if err = ms.All(
		func(g gattung.Gattung, m kennung.Matcher) (err error) {
			switch g {
			case gattung.Etikett:
				c.catEtiketten(u, m, aw)

			case gattung.Zettel:
				c.catZettelen(u, m, aw)

			case gattung.Typ:
				todo.Implement()
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c CatAlfred) catEtiketten(
	u *umwelt.Umwelt,
	m kennung.Matcher,
	aw *alfred.Writer,
) {
	var err error

	if err = u.StoreObjekten().GetKennungIndex().EachSchwanzen(
		func(e kennung_index.Indexed[kennung.Etikett]) (err error) {
			_, err = aw.WriteEtikett(e)
			return
		},
	); err != nil {
		aw.WriteError(err)
		return
	}
}

func (c CatAlfred) catZettelen(
	u *umwelt.Umwelt,
	m kennung.Matcher,
	aw *alfred.Writer,
) {
	var err error

	if err = u.StoreObjekten().Zettel().Query(
		m,
		aw.WriteZettelVerzeichnisse,
	); err != nil {
		aw.WriteError(err)
		return
	}

	return
}

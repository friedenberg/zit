package commands

import (
	"bufio"
	"flag"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/alfred"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type CatAlfred struct {
	Type gattung.Gattung
	Command
}

func init() {
	registerCommand(
		"cat-alfred",
		func(f *flag.FlagSet) Command {
			c := &CatAlfred{
				Type: gattung.Unknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithIds{
				CommandWithIds: c,
			}
		},
	)
}

func (c CatAlfred) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet(
		id_set.ProtoId{
			MutableId: &kennung.Konfig{},
		},
		id_set.ProtoId{
			MutableId: &hinweis.Hinweis{},
			Expand: func(v string) (out string, err error) {
				var h hinweis.Hinweis
				h, err = u.StoreObjekten().Abbr().ExpandHinweisString(v)
				out = h.String()
				return
			},
		},
		id_set.ProtoId{
			MutableId: &kennung.Etikett{},
			Expand: func(v string) (out string, err error) {
				var e kennung.Etikett
				e, err = u.StoreObjekten().Abbr().ExpandEtikettString(v)
				out = e.String()
				return
			},
		},
		id_set.ProtoId{
			MutableId: &kennung.Typ{},
		},
		id_set.ProtoId{
			MutableId: &ts.Time{},
		},
	)

	return
}

func (c CatAlfred) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	//this command does its own error handling
	wo := bufio.NewWriter(u.Out())
	defer errors.Deferred(&err, wo.Flush)

	var aw *alfred.Writer

	if aw, err = alfred.New(wo, u.StoreObjekten()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, aw.Close)

	wg := &sync.WaitGroup{}

	switch c.Type {
	case gattung.Etikett:
		wg.Add(1)
		go c.catEtiketten(u, ids, aw, wg)

	case gattung.Akte:
		wg.Add(1)
		go c.catZettelen(u, ids, aw, wg)

	case gattung.Zettel:
		wg.Add(1)
		go c.catZettelen(u, ids, aw, wg)

	case gattung.Hinweis:
		wg.Add(1)
		go c.catZettelen(u, ids, aw, wg)

	default:
		wg.Add(2)
		go c.catEtiketten(u, ids, aw, wg)
		go c.catZettelen(u, ids, aw, wg)
	}

	wg.Wait()

	return
}

func (c CatAlfred) catEtiketten(
	u *umwelt.Umwelt,
	ids id_set.Set,
	aw *alfred.Writer,
	wg *sync.WaitGroup,
) {
	if wg != nil {
		defer wg.Done()
	}

	var ea []kennung.Etikett

	var err error

	if ea, err = u.StoreObjekten().Etiketten(); err != nil {
		aw.WriteError(err)
		return
	}

	for _, e := range ea {
		aw.WriteEtikett(e)
	}
}

func (c CatAlfred) catZettelen(
	u *umwelt.Umwelt,
	ids id_set.Set,
	aw *alfred.Writer,
	wg *sync.WaitGroup,
) {
	if wg != nil {
		defer wg.Done()
	}

	wk := zettel.MakeWriterKonfig(u.Konfig())

	var err error

	if err = u.StoreObjekten().Zettel().ReadAllSchwanzen(
		collections.MakeChain(
			wk,
			zettel.WriterIds{
				Filter: id_set.Filter{
					Set:        ids,
					AllowEmpty: true,
				},
			}.WriteZettelTransacted,
			aw.WriteZettelVerzeichnisse,
		),
	); err != nil {
		aw.WriteError(err)
		return
	}

	return
}

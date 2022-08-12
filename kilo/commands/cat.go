package commands

import (
	"encoding/json"
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/bravo/zk_types"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type Cat struct {
	zk_types.Type
	Format string
}

func init() {
	registerCommand(
		"cat",
		func(f *flag.FlagSet) Command {
			c := &Cat{
				Type: zk_types.TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")
			f.StringVar(&c.Format, "format", "", "ObjekteType")

			return commandWithLockedStore{c}
		},
	)
}

func (c Cat) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	switch c.Type {
	case zk_types.TypeEtikett:
		err = c.etiketten(store)

	case zk_types.TypeZettel:
		err = c.zettelen(store)

	case zk_types.TypeAkte:
		err = c.akten(store)

	case zk_types.TypeHinweis:
		err = c.hinweisen(store)

	default:
		err = c.all(store)
	}

	return
}

func (c Cat) etiketten(store store_with_lock.Store) (err error) {
	var ea []etikett.Etikett

	if ea, err = store.Zettels().Etiketten(); err != nil {
		err = errors.Error(err)
		return
	}

OUTER:
	for _, e := range ea {
		prefixes := e.Expanded(etikett.ExpanderRight{})

	INNER:
		for tn, tv := range store.Konfig.Tags {
			if !tv.Hide {
				continue INNER
			}

			if prefixes.ContainsString(tn) {
				continue OUTER
			}
		}

		stdprinter.Out(e)
	}

	return
}

func (c Cat) zettelen(store store_with_lock.Store) (err error) {
	var all map[hinweis.Hinweis]stored_zettel.Transacted

	logz.Print()
	defer logz.Print()

	if all, err = store.Zettels().ZettelTails(); err != nil {
		err = errors.Error(err)
		return
	}

	if c.Format == "json" {

		// not a bottleneck
		for _, z := range all {
			b, err := json.Marshal(z.Stored)

			if err != nil {
				logz.Print(err)
			} else {
				stdprinter.Out(string(b))
			}
		}
	} else {
		f := zettel_formats.Text{}

		c := zettel.FormatContextWrite{
			Out: store.Out,
		}

		// not a bottleneck
		for _, z := range all {

			c.Zettel = z.Zettel

			if _, err = f.WriteTo(c); err != nil {
				logz.Print(err)
			}
		}
	}

	return
}

func (c Cat) akten(store store_with_lock.Store) (err error) {
	var shas []sha.Sha

	if shas, err = store.Akten().All(); err != nil {
		err = errors.Error(err)
		return
	}

	for _, s := range shas {
		stdprinter.Out(s)
	}

	return
}

func (c Cat) hinweisen(store store_with_lock.Store) (err error) {
	// var hins []hinweis.Hinweis
	// var shas []sha.Sha

	// if shas, hins, err = store.Hinweisen().All(); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	// for i, h := range hins {
	// 	stdprinter.Outf("%s: %s\n", h, shas[i])
	// }

	return
}

func (c Cat) all(store store_with_lock.Store) (err error) {
	// var hins []hinweis.Hinweis

	// if _, hins, err = store.Hinweisen().All(); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	// chains := make([]zettels.Chain, len(hins))

	// for i, h := range hins {
	// 	if chains[i], err = store.Zettels().AllInChain(h); err != nil {
	// 		err = errors.Error(err)
	// 		return
	// 	}
	// }

	// b, err := json.Marshal(chains)

	// if err != nil {
	// 	logz.Print(err)
	// } else {
	// 	stdprinter.Out(string(b))
	// }

	return
}

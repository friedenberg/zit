package commands

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type Cat struct {
	Type   _Type
	Format string
}

func init() {
	registerCommand(
		"cat",
		func(f *flag.FlagSet) Command {
			c := &Cat{
				Type: _TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")
			f.StringVar(&c.Format, "format", "", "ObjekteType")

			return commandWithLockedStore{c}
		},
	)
}

func (c Cat) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	switch c.Type {
	case _TypeEtikett:
		err = c.etiketten(store)

	case _TypeZettel:
		err = c.zettelen(store)

	case _TypeAkte:
		err = c.akten(store)

	case _TypeHinweis:
		err = c.hinweisen(store)

	default:
		err = c.all(store)
	}

	return
}

func (c Cat) etiketten(store store_with_lock.Store) (err error) {
	var ea []_Etikett

	if ea, err = store.Etiketten().All(); err != nil {
		err = errors.Error(err)
		return
	}

OUTER:
	for _, e := range ea {
		prefixes := e.Expanded(_EtikettExpanderRight{})

	INNER:
		for tn, tv := range store.Konfig.Tags {
			if !tv.Hide {
				continue INNER
			}

			if prefixes.ContainsString(tn) {
				continue OUTER
			}
		}

		_Out(e)
	}

	return
}

func (c Cat) zettelen(store store_with_lock.Store) (err error) {
	var all map[string]_NamedZettel

	if all, err = store.Zettels().All(); err != nil {
		err = errors.Error(err)
		return
	}

	if c.Format == "json" {

		// not a bottleneck
		for _, z := range all {
			b, err := json.Marshal(z.Stored)

			if err != nil {
				log.Print(err)
			} else {
				_Out(string(b))
			}
		}
	} else {
		f := _ZettelFormatsText{}

		c := _ZettelFormatContextWrite{
			Out: store.Out,
		}

		// not a bottleneck
		for _, z := range all {

			c.Zettel = z.Zettel

			if _, err = f.WriteTo(c); err != nil {
				log.Print(err)
			}
		}
	}

	return
}

func (c Cat) akten(store store_with_lock.Store) (err error) {
	var shas []_Sha

	if shas, err = store.Akten().All(); err != nil {
		err = errors.Error(err)
		return
	}

	for _, s := range shas {
		_Out(s)
	}

	return
}

func (c Cat) hinweisen(store store_with_lock.Store) (err error) {
	var hins []_Hinweis
	var shas []_Sha

	if shas, hins, err = store.Hinweisen().All(); err != nil {
		err = errors.Error(err)
		return
	}

	for i, h := range hins {
		_Outf("%s: %s\n", h, shas[i])
	}

	return
}

func (c Cat) all(store store_with_lock.Store) (err error) {
	var hins []_Hinweis

	if _, hins, err = store.Hinweisen().All(); err != nil {
		err = errors.Error(err)
		return
	}

	chains := make([]_ZettelsChain, len(hins))

	for i, h := range hins {
		if chains[i], err = store.Zettels().AllInChain(h); err != nil {
			err = errors.Error(err)
			return
		}
	}

	b, err := json.Marshal(chains)

	if err != nil {
		log.Print(err)
	} else {
		_Out(string(b))
	}

	return
}

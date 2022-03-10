package commands

import (
	"encoding/json"
	"flag"
	"log"
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

			return commandWithZettels{c}
		},
	)
}

func (c Cat) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	switch c.Type {
	case _TypeEtikett:
		err = c.etiketten(u, zs)

	case _TypeZettel:
		err = c.zettelen(u, zs)

	case _TypeAkte:
		err = c.akten(u, zs)

	case _TypeHinweis:
		err = c.hinweisen(u, zs)

	default:
		err = c.all(u, zs)
	}

	return
}

func (c Cat) etiketten(u _Umwelt, zs _Zettels) (err error) {
	var ea []_Etikett

	if ea, err = zs.Etiketten().All(); err != nil {
		err = _Error(err)
		return
	}

OUTER:
	for _, e := range ea {
		prefixes := e.Expanded(_EtikettExpanderRight{})

	INNER:
		for tn, tv := range u.Konfig.Tags {
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

func (c Cat) zettelen(u _Umwelt, zs _Zettels) (err error) {
	var all map[string]_NamedZettel

	if all, err = zs.All(); err != nil {
		err = _Error(err)
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
			Out: u.Out,
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

func (c Cat) akten(u _Umwelt, zs _Zettels) (err error) {
	var shas []_Sha

	if shas, err = zs.Akten().All(); err != nil {
		err = _Error(err)
		return
	}

	for _, s := range shas {
		_Out(s)
	}

	return
}

func (c Cat) hinweisen(u _Umwelt, zs _Zettels) (err error) {
	var hins []_Hinweis
	var shas []_Sha

	if shas, hins, err = zs.Hinweisen().All(); err != nil {
		err = _Error(err)
		return
	}

	for i, h := range hins {
		_Outf("%s: %s\n", h, shas[i])
	}

	return
}

func (c Cat) all(u _Umwelt, zs _Zettels) (err error) {
	var hins []_Hinweis

	if _, hins, err = zs.Hinweisen().All(); err != nil {
		err = _Error(err)
		return
	}

	chains := make([]_ZettelsChain, len(hins))

	for i, h := range hins {
		if chains[i], err = zs.AllInChain(h); err != nil {
			err = _Error(err)
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

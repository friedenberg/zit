package commands

import (
	"encoding/json"
	"flag"
	"sort"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/typ"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type Cat struct {
	gattung.Gattung
	Format string
}

func init() {
	registerCommand(
		"cat",
		func(f *flag.FlagSet) Command {
			c := &Cat{
				Gattung: gattung.Unknown,
			}

			f.Var(&c.Gattung, "gattung", "ObjekteType")
			f.StringVar(&c.Format, "format", "", "ObjekteType")

			return c
		},
	)
}

//TODO move to types as args
func (c Cat) Run(u *umwelt.Umwelt, args ...string) (err error) {
	switch c.Gattung {
	case gattung.Etikett:
		err = c.etiketten(u)

	case gattung.Zettel:
		err = c.zettelen(u)

	case gattung.Akte:
		err = c.akten(u)

	case gattung.Hinweis:
		err = c.hinweisen(u)

	case gattung.Typ:
		err = c.typen(u)

	default:
		err = c.all(u)
	}

	return
}

func (c Cat) etiketten(u *umwelt.Umwelt) (err error) {
	var ea []etikett.Etikett

	if ea, err = u.StoreObjekten().Etiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, e := range ea {
		if err = errors.PrintOut(e.String()); err != nil {
			err = errors.IsAsNilOrWrapf(
				err,
				syscall.EPIPE,
				"Etikett: %s",
				e,
			)

			return
		}
	}

	return
}

func (c Cat) zettelen(u *umwelt.Umwelt) (err error) {
	var all zettel_transacted.Set

	if all, err = u.StoreObjekten().ZettelenSchwanzen(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.Format == "json" {

		// not a bottleneck
		all.Each(
			func(z zettel_transacted.Zettel) (err error) {
				var b []byte

				b, err = json.Marshal(z.Named.Stored)

				if err != nil {
					err = errors.PrintErr(err)
				} else {
					err = errors.PrintOut(string(b))
				}

				if err != nil {
					//TODO combined error
					err = errors.IsAsNilOrWrapf(
						err,
						syscall.EPIPE,
						"Zettel: %s",
						z.Named.Hinweis,
					)

					return
				}

				return
			},
		)
	} else {
		f := zettel.Text{}

		c := zettel.FormatContextWrite{
			Out: u.Out(),
		}

		// not a bottleneck
		all.Each(
			func(z zettel_transacted.Zettel) (err error) {
				c.Zettel = z.Named.Stored.Zettel

				if _, err = f.WriteTo(c); err != nil {
					err = errors.IsAsNilOrWrapf(
						err,
						syscall.EPIPE,
						"Zettel: %s",
						z.Named.Hinweis,
					)

					return
				}

				return
			},
		)
	}

	return
}

func (c Cat) akten(u *umwelt.Umwelt) (err error) {
	var shas []sha.Sha

	if shas, err = u.StoreAkten().All(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, s := range shas {
		if err = errors.PrintOut(s); err != nil {
			if errors.Is(err, syscall.EPIPE) {
				err = nil
				break
			} else {
				errors.Print(err)
			}
		}
	}

	return
}

func (c Cat) hinweisen(u *umwelt.Umwelt) (err error) {
	return
}

func (c Cat) typen(u *umwelt.Umwelt) (err error) {
	var all zettel_transacted.Set

	if all, err = u.StoreObjekten().ZettelenSchwanzen(); err != nil {
		err = errors.Wrap(err)
		return
	}

	typen := make(map[typ.Typ]bool)

	all.Each(
		func(z zettel_transacted.Zettel) (err error) {
			typen[z.Named.Stored.Zettel.Typ] = true

			return
		},
	)

	sortedTypen := make([]typ.Typ, 0, len(typen))

	for t, _ := range typen {
		sortedTypen = append(sortedTypen, t)
	}

	sort.Slice(sortedTypen, func(i, j int) bool { return sortedTypen[i].Less(sortedTypen[j].Etikett) })

	for _, t := range sortedTypen {
		errors.PrintOut(t)
	}

	return
}

func (c Cat) all(u *umwelt.Umwelt) (err error) {
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

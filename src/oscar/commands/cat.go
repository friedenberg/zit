package commands

import (
	"flag"
	"fmt"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/typ"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
	"github.com/friedenberg/zit/src/mike/umwelt"
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

			return commandWithIds{c}
		},
	)
}

// TODO move to types as args
func (c Cat) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	switch c.Gattung {
	case gattung.Etikett:
		err = c.etiketten(u)

	case gattung.Zettel:
		err = c.zettelen(u)

	case gattung.Akte:
		err = c.akten(u, ids)

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

func (c Cat) zettelWriter(
	u *umwelt.Umwelt,
) collections.WriterFunc[*zettel.Zettel] {
	switch c.Format {
	case "json":
		return zettel.MakeSerializedFormatWriter(
			zettel.JsonObjekte{},
			u.Out(),
			u.StoreObjekten(),
			u.Konfig(),
		)

	default:
		return zettel.MakeSerializedFormatWriter(
			zettel.Text{},
			u.Out(),
			u.StoreObjekten(),
			u.Konfig(),
		)
	}
}

func (c Cat) zettelen(u *umwelt.Umwelt) (err error) {
	w := zettel_transacted.MakeWriterChain(
		zettel_transacted.MakeWriterZettel(
			c.zettelWriter(u),
		),
	)

	if err = u.StoreObjekten().ReadAllSchwanzenTransacted(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Cat) akten(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	if err = u.StoreObjekten().ReadAllAktenShas(
		func(s sha.Sha) (err error) {
			_, err = fmt.Fprintln(u.Out(), s)
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}
	// w := zettel_transacted.MakeWriterChain(
	// 	zettel_transacted.MakeWriterZettelNamed(
	// 		zettel_named.FilterIdSet{
	// 			AllowEmpty: true,
	// 			Set:        ids,
	// 		}.WriteZettelNamed,
	// 	),
	// 	zettel_transacted.MakeWriterZettel(
	// 		collections.MakeSyncSerializer(
	// 			func(z *zettel.Zettel) (err error) {
	// 				_, err = fmt.Fprintln(u.Out(), z.Akte.String())

	// 				return
	// 			},
	// 		),
	// 	),
	// )

	// if err = u.StoreObjekten().ReadAllSchwanzenTransacted(w); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	return
}

func (c Cat) hinweisen(u *umwelt.Umwelt) (err error) {
	return
}

func (c Cat) typen(u *umwelt.Umwelt) (err error) {
	typen := collections.MakeMutableValueSet[typ.Typ, *typ.Typ]()

	if err = u.StoreObjekten().ReadAllSchwanzenTransacted(
		zettel_transacted.MakeWriter(
			func(z *zettel_transacted.Zettel) (err error) {
				err = typen.Add(z.Named.Stored.Zettel.Typ)

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sortedTypen := typen.Copy().SortedString()

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

package commands

import (
	"flag"
	"fmt"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/id_set"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/november/umwelt"
)

// TODO-P3 examine removing cat entirely
type Cat struct {
	gattung.Gattung
	//Specific to Gattung
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

// TODO-P3 move to types as args
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
	var ea []kennung.Etikett

	if ea, err = u.StoreObjekten().Etiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, e := range ea {
		if err = errors.Out().Print(e.String()); err != nil {
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
) collections.WriterFunc[*zettel.Objekte] {
	switch c.Format {
	case "json":
		return zettel.MakeSerializedFormatWriter(
			zettel.MakeObjekteFormatterJson(),
			u.Out(),
			u.StoreObjekten(),
			u.Konfig(),
		)

	default:
		return zettel.MakeSerializedFormatWriter(
			zettel.MakeObjekteTextFormatterIncludeAkte(
				u.Standort(),
				u.Konfig(),
				u.StoreObjekten(),
				nil,
			),
			u.Out(),
			u.StoreObjekten(),
			u.Konfig(),
		)
	}
}

func (c Cat) zettelen(u *umwelt.Umwelt) (err error) {
	if err = u.StoreObjekten().Zettel().ReadAllSchwanzen(
		zettel.MakeWriterZettel(
			c.zettelWriter(u),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

//func (c CatObjekte) akten(store *umwelt.Umwelt, ids id_set.Set) (err error) {
//	type akteToWrite struct {
//		io.ReadCloser
//		*zettel_named.Zettel
//	}

//	akteWriter := collections.MakeSyncSerializer(
//		func(a akteToWrite) (err error) {
//			errors.Log().Printf("writing one: %s", a)
//			defer errors.Deferred(&err, a.ReadCloser.Close)

//			//TODO-P2 explicitly support toml
//			if _, err = io.WriteString(
//				store.Out(),
//				fmt.Sprintf("['%s']\n", a.Hinweis),
//			); err != nil {
//				err = errors.Wrap(err)
//				return
//			}

//			if _, err = io.Copy(store.Out(), a.ReadCloser); err != nil {
//				err = errors.Wrap(err)
//				return
//			}

//			return
//		},
//	)

//	if err = c.akteShasFromIds(
//		store,
//		ids,
//		func(z *zettel.Zettel) (err error) {
//			sb := z.Objekte.Akte

//			if sb.IsNull() {
//				return
//			}

//			var r io.ReadCloser

//			if r, err = store.StoreObjekten().AkteReader(sb); err != nil {
//				err = errors.Wrap(err)
//				return
//			}

//			if err = akteWriter(
//				akteToWrite{
//					ReadCloser: r,
//					Zettel:     &z.Named,
//				},
//			); err != nil {
//				err = errors.Wrap(err)
//				return
//			}

//			return
//		},
//	); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	return
//}

func (c Cat) akten(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	if err = u.Standort().ReadAllShasForGattung(
		gattung.Akte,
		func(s sha.Sha) (err error) {
			_, err = fmt.Fprintln(u.Out(), s)
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}
	// w := zettel.MakeWriterChain(
	// 	zettel.MakeWriterZettelNamed(
	// 		zettel_named.FilterIdSet{
	// 			AllowEmpty: true,
	// 			Set:        ids,
	// 		}.WriteZettelNamed,
	// 	),
	// 	zettel.MakeWriterZettel(
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
	typen := collections.MakeMutableValueSet[kennung.Typ, *kennung.Typ]()

	if err = u.StoreObjekten().Zettel().ReadAllSchwanzen(
		func(z *zettel.Transacted) (err error) {
			err = typen.Add(z.Objekte.Typ)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sortedTypen := typen.Copy().SortedString()

	for _, t := range sortedTypen {
		errors.Out().Print(t)
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

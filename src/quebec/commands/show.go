package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Show struct {
	Format string
}

func init() {
	registerCommand(
		"show",
		func(f *flag.FlagSet) Command {
			c := &Show{}

			f.StringVar(&c.Format, "format", "text", "format")

			return commandWithQuery{
				CommandWithQuery: c,
			}
		},
	)
}

func (c Show) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
		gattung.Bestandsaufnahme,
		gattung.Kasten,
	)
}

func (c Show) RunWithQuery(u *umwelt.Umwelt, ms kennung.MetaSet) (err error) {
	if err = ms.All(
		func(g gattung.Gattung, ids kennung.Set) (err error) {
			switch g {

			case gattung.Akte:
				return c.showAkten(u, ids)

			case gattung.Zettel:
				var fv zettel.FormatterValue

				if err = fv.Set(c.Format); err != nil {
					err = errors.Normal(err)
					return
				}

				return c.showOneOrMoreZettels(
					u,
					ids,
					fv.FuncFormatterVerzeichnisse(
						u.Out(),
						u.StoreObjekten(),
						u.Konfig(),
						u.PrinterZettelTransacted(),
					),
				)

			case gattung.Transaktion:
				return c.showTransaktions(u, ids)

			case gattung.Typ:
				var ev typ.FormatterValue

				if err = ev.Set(c.Format); err != nil {
					err = errors.Normal(err)
					return
				}

				return c.showTypen(
					u,
					ids,
					ev.FuncFormatter(
						u.Out(),
						u.StoreObjekten(),
						u.PrinterTypTransacted(),
					),
				)

			case gattung.Etikett:
				var ev etikett.FormatterValue

				if err = ev.Set(c.Format); err != nil {
					err = errors.Normal(err)
					return
				}

				return c.showEtiketten(
					u,
					ids,
					ev.FuncFormatter(
						u.Out(),
						u.StoreObjekten(),
					),
				)

			case gattung.Konfig:
				var ev erworben.FormatterValue

				if err = ev.Set(c.Format); err != nil {
					err = errors.Normal(err)
					return
				}

				return c.showKonfig(
					u,
					ev.FuncFormatter(
						u.Out(),
						u.StoreObjekten(),
					),
				)

			default:
				err = errors.Errorf("unsupported Gattung: %s", g)
				return
			}
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Show) showOneOrMoreZettels(
	u *umwelt.Umwelt,
	ids kennung.Set,
	fv schnittstellen.FuncIter[*zettel.Transacted],
) (err error) {
	if h, ok := ids.OnlySingleHinweis(); ok {
		if err = c.showOneZettel(u, h, fv); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = c.showManyZettels(u, ids, fv); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Show) showOneZettel(
	u *umwelt.Umwelt,
	h kennung.Hinweis,
	fv schnittstellen.FuncIter[*zettel.Transacted],
) (err error) {
	var z *zettel.Transacted

	if z, err = u.StoreWorkingDirectory().ReadOne(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fv(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Show) showManyZettels(
	u *umwelt.Umwelt,
	ids kennung.Set,
	fv schnittstellen.FuncIter[*zettel.Transacted],
) (err error) {
	idFilter := zettel.WriterIds{
		Filter: kennung.Filter{
			Set: ids,
		},
	}.WriteZettelTransacted

	method := u.StoreWorkingDirectory().ReadMany

	filter := idFilter

	if ids.Sigil.IncludesHistory() {
		method = u.StoreWorkingDirectory().ReadManyHistory
		hinweisen := kennung.MakeHinweisMutableSet()

		if err = u.StoreObjekten().Zettel().ReadAllSchwanzen(
			iter.MakeChain(
				idFilter,
				func(o *zettel.Transacted) (err error) {
					return hinweisen.Add(o.Sku.Kennung)
				},
			),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		hContainer := collections.WriterContainer[kennung.Hinweis](
			hinweisen,
			collections.MakeErrStopIteration(),
		)

		filter = func(o *zettel.Transacted) (err error) {
			err = hContainer(o.Sku.Kennung)

			if collections.IsStopIteration(err) {
				err = idFilter(o)
			}

			return
		}
	}

	f1 := collections.MakeSyncSerializer(fv)

	if err = method(
		iter.MakeChain(
			filter,
			f1,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO-P3 support All
func (c Show) showAkten(u *umwelt.Umwelt, ids kennung.Set) (err error) {
	zettels := make([]*zettel.Transacted, ids.Len())

	for i, is := range ids.AnyShasOrHinweisen() {
		var tz *zettel.Transacted

		if tz, err = u.StoreObjekten().Zettel().ReadOne(is); err != nil {
			err = errors.Wrap(err)
			return
		}

		zettels[i] = tz
	}

	var ar io.ReadCloser

	for _, named := range zettels {
		if ar, err = u.StoreObjekten().AkteReader(named.Objekte.Akte); err != nil {
			err = errors.Wrap(err)
			return
		}

		if ar == nil {
			err = errors.Errorf("akte reader is nil")
			return
		}

		defer errors.Deferred(&err, ar.Close)

		if _, err = io.Copy(u.Out(), ar); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

// TODO-P3 support All
func (c Show) showTransaktions(u *umwelt.Umwelt, ids kennung.Set) (err error) {
	ids.Timestamps.ImmutableClone().Each(
		func(is ts.Time) (err error) {
			var t *transaktion.Transaktion

			if t, err = u.StoreObjekten().GetTransaktionStore().ReadTransaktion(
				is,
			); err != nil {
				errors.PrintErrf("error: %s", err)
				return
			}

			errors.Out().Printf("transaktion: %#v", t)

			t.Skus.Each(
				func(o sku.SkuLike) (err error) {
					errors.Out().Print(o)
					return
				},
			)

			return
		},
	)

	return
}

func (c Show) showTypen(
	u *umwelt.Umwelt,
	ids kennung.Set,
	f schnittstellen.FuncIter[*typ.Transacted],
) (err error) {
	f1 := collections.MakeSyncSerializer(f)

	typen := ids.Typen.MutableClone()

	method := u.StoreObjekten().Typ().ReadAllSchwanzen

	if ids.Sigil.IncludesHistory() {
		method = u.StoreObjekten().Typ().ReadAll
	}

	if err = method(
		func(t *typ.Transacted) (err error) {
			switch {
			case ids.Sigil.IncludesSchwanzen():
				fallthrough

			case typen.Contains(t.Sku.Kennung):
				if err = f1(t); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO-P3 support All
func (c Show) showEtiketten(
	u *umwelt.Umwelt,
	ids kennung.Set,
	f schnittstellen.FuncIter[*etikett.Transacted],
) (err error) {
	f1 := collections.MakeSyncSerializer(f)

	etiketten := ids.Etiketten.MutableClone()
	if err = etiketten.EachPtr(
		iter.MakeChain(
			func(t *kennung.Etikett) (err error) {
				ty := u.Konfig().GetEtikett(*t)

				if ty == nil {
					return
				}

				if err = f1(ty); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Show) showKonfig(
	u *umwelt.Umwelt,
	f schnittstellen.FuncIter[*erworben.Transacted],
) (err error) {
	f1 := collections.MakeSyncSerializer(f)

	var k *erworben.Transacted

	if k, err = u.StoreObjekten().Konfig().Read(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = f1(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

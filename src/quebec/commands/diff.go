package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
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

type Diff struct {
	GattungSet gattungen.MutableSet
	Format     string
}

func init() {
	registerCommand(
		"diff",
		func(f *flag.FlagSet) Command {
			c := &Diff{
				GattungSet: gattungen.MakeMutableSet(gattung.Zettel),
			}

			gsvs := collections.MutableValueSet2[gattung.Gattung, *gattung.Gattung]{
				MutableSetLike: &c.GattungSet,
				SetterPolicy:   collections.SetterPolicyReset,
			}

			f.Var(gsvs, "gattung", "Gattung")
			f.StringVar(&c.Format, "format", "text", "format")

			cwi := commandWithIds{
				CommandWithIds: c,
			}

			return CommandV2{
				Command:        cwi,
				WithCompletion: cwi,
			}
		},
	)
}

func (c Diff) ProtoIdSet(u *umwelt.Umwelt) (is kennung.ProtoIdSet) {
	is = kennung.MakeProtoIdSet()

	if c.GattungSet.Contains(gattung.Zettel) {
		is.AddMany(
			kennung.ProtoId{
				Setter: kennung.MakeSigil(kennung.SigilNone),
			},
			// id_set.ProtoId{
			// 	Setter: &sha.Sha{},
			// },
			kennung.ProtoId{
				Setter: &kennung.Hinweis{},
				Expand: func(v string) (out string, err error) {
					var h kennung.Hinweis
					h, err = u.StoreObjekten().GetAbbrStore().ExpandHinweisString(v)
					out = h.String()
					return
				},
			},
			kennung.ProtoId{
				Setter: &kennung.Etikett{},
				Expand: func(v string) (out string, err error) {
					var e kennung.Etikett
					e, err = u.StoreObjekten().GetAbbrStore().ExpandEtikettString(v)
					out = e.String()
					return
				},
			},
			kennung.ProtoId{
				Setter: &kennung.Typ{},
			},
			// id_set.ProtoId{
			// 	Setter: &ts.Time{},
			// },
		)
	}

	if c.GattungSet.Contains(gattung.Typ) {
		is.AddMany(
			kennung.ProtoId{
				Setter: &kennung.Typ{},
			},
		)
	}

	if c.GattungSet.Contains(gattung.Transaktion) {
		is.AddMany(
			kennung.ProtoId{
				Setter: &ts.Time{},
			},
		)
	}

	return
}

func (c Diff) RunWithIds(u *umwelt.Umwelt, ids kennung.Set) (err error) {
	if err = c.GattungSet.Each(
		func(g gattung.Gattung) (err error) {
			switch g {

			// case gattung.Akte:
			// 	return c.showAkten(u, ids)

			case gattung.Zettel:
				var fv zettel.FormatterValue

				if err = fv.Set(c.Format); err != nil {
					err = errors.Normal(err)
					return
				}

				return c.showZettels(
					u,
					ids,
					fv.FuncFormatterVerzeichnisse(
						u.Out(),
						u.StoreObjekten(),
						u.Konfig(),
						u.PrinterZettelTransacted(),
					),
				)

			// case gattung.Transaktion:
			// 	return c.showTransaktions(u, ids)

			// case gattung.Typ:
			// 	var ev typ.FormatterValue

			// 	if err = ev.Set(c.Format); err != nil {
			// 		err = errors.Normal(err)
			// 		return
			// 	}

			// 	return c.showTypen(
			// 		u,
			// 		ids,
			// 		ev.FuncFormatter(
			// 			u.Out(),
			// 			u.StoreObjekten(),
			// 			u.PrinterTypTransacted(format.StringNew),
			// 		),
			// 	)

			// case gattung.Etikett:
			// 	var ev etikett.FormatterValue

			// 	if err = ev.Set(c.Format); err != nil {
			// 		err = errors.Normal(err)
			// 		return
			// 	}

			// 	return c.showEtiketten(
			// 		u,
			// 		ids,
			// 		ev.FuncFormatter(
			// 			u.Out(),
			// 			u.StoreObjekten(),
			// 		),
			// 	)

			// case gattung.Konfig:
			// 	var ev erworben.FormatterValue

			// 	if err = ev.Set(c.Format); err != nil {
			// 		err = errors.Normal(err)
			// 		return
			// 	}

			// 	return c.showKonfig(
			// 		u,
			// 		ev.FuncFormatter(
			// 			u.Out(),
			// 			u.StoreObjekten(),
			// 		),
			// 	)

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

func (c Diff) showZettels(
	u *umwelt.Umwelt,
	ids kennung.Set,
	fv collections.WriterFunc[*zettel.Transacted],
) (err error) {
	idFilter := zettel.WriterIds{
		Filter: kennung.Filter{
			Set: ids,
		},
	}.WriteZettelTransacted

	method := u.StoreWorkingDirectory().ReadMany
	// method := u.StoreObjekten().Zettel().ReadAll

	filter := idFilter

	if ids.Sigil.IncludesHistory() {
		method = u.StoreWorkingDirectory().ReadManyHistory
		// method = u.StoreObjekten().Zettel().ReadAll
		hinweisen := kennung.MakeHinweisMutableSet()

		if err = u.StoreObjekten().Zettel().ReadAllSchwanzen(
			collections.MakeChain(
				idFilter,
				func(o *zettel.Transacted) (err error) {
					return hinweisen.Add(o.Sku.Kennung)
				},
			),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		hContainer := hinweisen.WriterContainer(collections.MakeErrStopIteration())

		filter = func(o *zettel.Transacted) (err error) {
			err = hContainer(o.Sku.Kennung)

			if collections.IsStopIteration(err) {
				err = idFilter(o)
			}

			return
		}
	}

	f1 := collections.MakeSyncSerializer(fv)
	errors.Log().Printf("%v", filter)

	if err = method(
		collections.MakeChain(
			func(z *zettel.Transacted) (err error) {
				errors.Log().Printf("processing zettel: %s", z.Sku)
				return
			},
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
func (c Diff) showAkten(u *umwelt.Umwelt, ids kennung.Set) (err error) {
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
func (c Diff) showTransaktions(u *umwelt.Umwelt, ids kennung.Set) (err error) {
	ids.Timestamps.Copy().Each(
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

func (c Diff) showTypen(
	u *umwelt.Umwelt,
	ids kennung.Set,
	f collections.WriterFunc[*typ.Transacted],
) (err error) {
	f1 := collections.MakeSyncSerializer(f)

	typen := ids.Typen.MutableCopy()

	method := u.StoreObjekten().Typ().ReadAllSchwanzen

	if ids.Sigil.IncludesHistory() {
		method = u.StoreObjekten().Typ().ReadAll
	}

	if err = method(
		func(t *typ.Transacted) (err error) {
			switch {
			// case c.All:
			// 	fallthrough

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
func (c Diff) showEtiketten(
	u *umwelt.Umwelt,
	ids kennung.Set,
	f collections.WriterFunc[*etikett.Transacted],
) (err error) {
	f1 := collections.MakeSyncSerializer(f)

	etiketten := ids.Etiketten.Copy().MutableCopy()
	if err = etiketten.EachPtr(
		collections.MakeChain(
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

func (c Diff) showKonfig(
	u *umwelt.Umwelt,
	f collections.WriterFunc[*erworben.Transacted],
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
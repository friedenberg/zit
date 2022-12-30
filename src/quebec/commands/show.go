package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/india/etikett"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/india/typ"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Show struct {
	//TODO-P0 gattung.Set
	gattung.Gattung
	Format string
	All    bool
}

func init() {
	registerCommand(
		"show",
		func(f *flag.FlagSet) Command {
			c := &Show{
				Gattung: gattung.Zettel,
			}

			f.Var(&c.Gattung, "gattung", "Gattung")
			f.StringVar(&c.Format, "format", "text", "format")
			f.BoolVar(&c.All, "all", false, "show all Objekten")

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

func (c Show) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	switch c.Gattung {

	default:
		is = id_set.MakeProtoIdSet(
			id_set.ProtoId{
				MutableId: &sha.Sha{},
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

	case gattung.Typ:
		is = id_set.MakeProtoIdSet(
			id_set.ProtoId{
				MutableId: &kennung.Typ{},
			},
		)

	case gattung.Transaktion:
		is = id_set.MakeProtoIdSet(
			id_set.ProtoId{
				MutableId: &ts.Time{},
			},
		)
	}

	return
}

func (c Show) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	switch c.Gattung {

	case gattung.Akte:
		return c.showAkten(u, ids)

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
		var ev konfig.FormatterValue

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
		err = errors.Errorf("unsupported Gattung: %s", c.Gattung)
		return
	}
}

func (c Show) showZettels(
	u *umwelt.Umwelt,
	ids id_set.Set,
	fv collections.WriterFunc[*zettel.Transacted],
) (err error) {
	filter := zettel.WriterIds{
		Filter: id_set.Filter{
			AllowEmpty: c.All,
			Set:        ids,
		},
	}.WriteZettelVerzeichnisse

	method := u.StoreWorkingDirectory().ReadMany

	if u.Konfig().IncludeHistory {
		method = u.StoreWorkingDirectory().ReadManyHistory
		hinweisen := hinweis.MakeMutableSet()

		if err = u.StoreObjekten().Zettel().ReadAllSchwanzenVerzeichnisse(
			collections.MakeChain(
				filter,
				func(o *zettel.Transacted) (err error) {
					return hinweisen.Add(o.Sku.Kennung)
				},
			),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		hContainer := hinweisen.WriterContainer(io.EOF)

		filter = func(o *zettel.Transacted) (err error) {
			return hContainer(o.Sku.Kennung)
		}
	}

	f1 := collections.MakeSyncSerializer(fv)

	if err = method(
		collections.MakeChain(
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
func (c Show) showAkten(u *umwelt.Umwelt, ids id_set.Set) (err error) {
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
func (c Show) showTransaktions(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	ids.Timestamps.Copy().Each(
		func(is ts.Time) (err error) {
			var t *transaktion.Transaktion

			if t, err = u.StoreObjekten().ReadTransaktion(is); err != nil {
				errors.PrintErrf("error: %s", err)
				return
			}

			errors.Out().Printf("transaktion: %#v", t)

			t.Each(
				func(o *sku.Sku) (err error) {
					errors.Out().Print(o)
					return
				},
			)

			return
		},
	)

	return
}

// TODO-P3 support All
func (c Show) showTypen(
	u *umwelt.Umwelt,
	ids id_set.Set,
	f collections.WriterFunc[*typ.Transacted],
) (err error) {
	f1 := collections.MakeSyncSerializer(f)

	typen := ids.Typen.MutableCopy()

	if err = typen.EachPtr(
		collections.MakeChain(
			func(t *kennung.Typ) (err error) {
				ty := u.Konfig().GetTyp(*t)

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

// TODO-P3 support All
func (c Show) showEtiketten(
	u *umwelt.Umwelt,
	ids id_set.Set,
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

func (c Show) showKonfig(
	u *umwelt.Umwelt,
	f collections.WriterFunc[*konfig.Transacted],
) (err error) {
	f1 := collections.MakeSyncSerializer(f)

	var k *konfig.Transacted

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

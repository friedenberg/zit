package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/kilo/objekte_formatter"
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

func (c Show) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
		// gattung.Bestandsaufnahme,
		gattung.Kasten,
	)
}

func (c Show) runGenericObjekteFormatterValue(
	u *umwelt.Umwelt,
	ms kennung.MetaSet,
	objekteFormatterValue objekte.FormatterValue,
) (err error) {
	f := collections.MakeSyncSerializer(
		objekteFormatterValue.MakeFormatterObjekte(
			u.Out(),
			u.StoreObjekten(),
			u.Konfig(),
			u.PrinterTransactedLike(),
		),
	)

	if err = u.StoreObjekten().Query(ms, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Show) RunWithQuery(u *umwelt.Umwelt, ms kennung.MetaSet) (err error) {
	objekteFormatterValue := objekte.FormatterValue{}

	if err = objekteFormatterValue.Set(c.Format); err == nil {
		return c.runGenericObjekteFormatterValue(u, ms, objekteFormatterValue)
	}

	err = nil

	var f objekte_formatter.Formatter

	if f, err = objekte_formatter.MakeFormatter(
		ms,
		c.Format,
		u.Out(),
		u.StoreObjekten(),
		u.Konfig(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.StoreObjekten().Query(
		ms,
		collections.MakeSyncSerializer(
			f.MakeFormatFunc(),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// func (c Show) showOneOrMoreZettels(
// 	u *umwelt.Umwelt,
// 	m kennung.Matcher,
// 	fv schnittstellen.FuncIter[*zettel.Transacted],
// ) (err error) {
// 	// errors.Err().Printf("ids: %s", ids)
// 	if h, ok := ids.OnlySingleHinweis(); ok {
// 		// errors.Err().Printf("only one")
// 		if err = c.showOneZettel(u, h, ids.Sigil, fv); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	} else {
// 		// errors.Err().Printf("many")
// 		if err = c.showManyZettels(u, ids, fv); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	return
// }

// func (c Show) showOneZettel(
// 	u *umwelt.Umwelt,
// 	h kennung.Hinweis,
// 	s kennung.Sigil,
// 	fv schnittstellen.FuncIter[*zettel.Transacted],
// ) (err error) {
// 	var z *zettel.Transacted

// 	if s.IncludesCwd() {
// 		if z, err = u.StoreWorkingDirectory().ReadOne(h); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	} else {
// 		if z, err = u.StoreObjekten().Zettel().ReadOne(h); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	if err = fv(z); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// func (c Show) showManyZettels(
// 	u *umwelt.Umwelt,
// 	ids kennung.Set,
// 	fv schnittstellen.FuncIter[*zettel.Transacted],
// ) (err error) {
// 	f1 := collections.MakeSyncSerializer(fv)

// 	if err = u.StoreObjekten().Zettel().Query(
// 		ids,
// 		f1,
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// func (c Show) showTypen(
// 	u *umwelt.Umwelt,
// 	m kennung.Matcher,
// 	f schnittstellen.FuncIter[*typ.Transacted],
// ) (err error) {
// 	f1 := collections.MakeSyncSerializer(f)

// 	typen := ids.Typen.MutableClone()

// 	method := u.StoreObjekten().Typ().ReadAllSchwanzen

// 	if ids.Sigil.IncludesHistory() {
// 		method = u.StoreObjekten().Typ().ReadAll
// 	}

// 	if err = method(
// 		func(t *typ.Transacted) (err error) {
// 			switch {
// 			case ids.Sigil.IncludesSchwanzen():
// 				fallthrough

// 			case typen.Contains(t.Sku.Kennung):
// 				if err = f1(t); err != nil {
// 					err = errors.Wrap(err)
// 					return
// 				}
// 			}

// 			return
// 		},
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// // TODO-P3 support All
// func (c Show) showEtiketten(
// 	u *umwelt.Umwelt,
// 	ids kennung.Set,
// 	f schnittstellen.FuncIter[*etikett.Transacted],
// ) (err error) {
// 	f1 := collections.MakeSyncSerializer(f)

// 	etiketten := ids.Etiketten.MutableClone()
// 	if err = etiketten.GetIncludes().EachPtr(
// 		iter.MakeChain(
// 			func(t *kennung.Etikett) (err error) {
// 				ty := u.Konfig().GetEtikett(*t)

// 				if ty == nil {
// 					return
// 				}

// 				if err = f1(ty); err != nil {
// 					err = errors.Wrap(err)
// 					return
// 				}

// 				return
// 			},
// 		),
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// func (c Show) showKonfig(
// 	u *umwelt.Umwelt,
// 	f schnittstellen.FuncIter[*erworben.Transacted],
// ) (err error) {
// 	f1 := collections.MakeSyncSerializer(f)

// 	var k *erworben.Transacted

// 	if k, err = u.StoreObjekten().Konfig().Read(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if err = f1(k); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

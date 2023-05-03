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
	registerCommandWithQuery(
		"show",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Show{}

			f.StringVar(&c.Format, "format", "text", "format")

			return c
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

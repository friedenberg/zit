package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/objekte_formatter"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Show struct {
	Format string
}

func init() {
	registerCommandWithQuery(
		"show",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Show{}

			f.StringVar(&c.Format, "format", "log", "format")

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
	ms matcher.Query,
	objekteFormatterValue objekte.FormatterValue,
) (err error) {
	f := iter.MakeSyncSerializer(
		objekteFormatterValue.MakeFormatterObjekte(
			u.Out(),
			u.Standort(),
			u.Konfig(),
			u.PrinterTransactedLike(),
			u.StringFormatWriterSkuTransactedShort(),
			u.StoreUtil().GetEnnui(),
			u.StoreUtil().ReadOneEnnui,
		),
	)

	// f := func(z *sku.Transacted) (err error) {
	// 	var shas []*sha.Sha

	// 	if shas, err = u.StoreUtil().GetEnnui().Get(z.GetMetadatei()); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	log.Debug().Printf("%s", shas)

	// 	return
	// }

	if err = u.StoreObjekten().QueryWithCwd(ms, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Show) RunWithQuery(u *umwelt.Umwelt, ms matcher.Query) (err error) {
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
		u.Standort(),
		u.Konfig(),
		u.StoreObjekten().GetAkten().GetTypV0(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.StoreObjekten().QueryWithCwd(
		ms,
		iter.MakeSyncSerializer(f.MakeFormatFunc()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
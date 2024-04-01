package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/juliett/objekte"
	"code.linenisgreat.com/zit/src/kilo/objekte_formatter"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
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

func (c Show) CompletionGattung() kennung.Gattung {
	return kennung.MakeGattung(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
		gattung.Bestandsaufnahme,
		gattung.Kasten,
	)
}

func (c Show) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
		// gattung.Bestandsaufnahme,
		gattung.Kasten,
	)
}

func (c Show) runGenericObjekteFormatterValue(
	u *umwelt.Umwelt,
	ms *query.Group,
	objekteFormatterValue objekte.FormatterValue,
) (err error) {
	f := iter.MakeSyncSerializer(
		objekteFormatterValue.MakeFormatterObjekte(
			u.Out(),
			u.Standort(),
			u.Konfig(),
			u.PrinterTransactedLike(),
			u.StringFormatWriterSkuTransactedShort(),
			u.GetStore().GetEnnui(),
			u.GetStore().ReadOneEnnui,
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

	if err = u.GetStore().QueryWithCwd(ms, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Show) RunWithQuery(u *umwelt.Umwelt, ms *query.Group) (err error) {
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
		u.GetStore().GetAkten().GetTypV0(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetStore().QueryWithCwd(
		ms,
		iter.MakeSyncSerializer(f.MakeFormatFunc()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

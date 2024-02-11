package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/delta/gattungen"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/objekte"
	"code.linenisgreat.com/zit/src/lima/bestandsaufnahme"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type Last struct {
	Type   gattung.Gattung
	Format string
}

func init() {
	registerCommand(
		"last",
		func(f *flag.FlagSet) Command {
			c := &Last{}

			f.StringVar(&c.Format, "format", "log", "format")

			return c
		},
	)
}

func (c Last) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Bestandsaufnahme,
	)
}

func (c Last) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) != 0 {
		errors.Err().Print("ignoring arguments")
	}

	var f schnittstellen.FuncIter[*sku.Transacted]

	objekteFormatterValue := objekte.FormatterValue{}

	if err = objekteFormatterValue.Set(c.Format); err != nil {
		err = errors.Wrap(err)
		return
	}

	f = objekteFormatterValue.MakeFormatterObjekte(
		u.Out(),
		u.Standort(),
		u.Konfig(),
		u.PrinterTransactedLike(),
		u.StringFormatWriterSkuTransactedShort(),
		u.StoreUtil().GetEnnui(),
		u.StoreUtil().ReadOneEnnui,
	)

	if err = c.runWithBestandsaufnahm(u, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Last) runWithBestandsaufnahm(
	u *umwelt.Umwelt,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	s := u.StoreObjekten()

	var b *sku.Transacted

	if b, err = s.GetBestandsaufnahmeStore().ReadLast(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var a *bestandsaufnahme.Akte

	if a, err = s.GetBestandsaufnahmeStore().GetAkte(b.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.TodoP3("support log line format for skus")
	if err = a.Skus.EachPtr(
		func(sk *sku.Transacted) (err error) {
			return f(sk)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

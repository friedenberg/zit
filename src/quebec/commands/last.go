package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/bestandsaufnahme"
	"github.com/friedenberg/zit/src/oscar/umwelt"
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

	var f schnittstellen.FuncIter[sku.SkuLikePtr]

	objekteFormatterValue := objekte.FormatterValue{}

	if err = objekteFormatterValue.Set(c.Format); err != nil {
		err = errors.Wrap(err)
		return
	}

	f = objekteFormatterValue.MakeFormatterObjekte(
		u.Out(),
		u.StoreObjekten(),
		u.Konfig(),
		u.PrinterTransactedLike(),
		u.StringFormatWriterSkuLikePtrShort(),
	)

	if err = c.runWithBestandsaufnahm(u, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Last) runWithBestandsaufnahm(
	u *umwelt.Umwelt,
	f schnittstellen.FuncIter[sku.SkuLikePtr],
) (err error) {
	s := u.StoreObjekten()

	var b *bestandsaufnahme.Transacted

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
	if err = sku.HeapEachPtr(a.Skus, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

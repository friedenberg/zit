package kennung_fmt

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

func MakeAligned(
	options erworben_cli_print_options.PrintOptions,
	ab kennung.Abbr,
) *aligned {
	return &aligned{PrintOptions: options, Abbr: ab}
}

type aligned struct {
	erworben_cli_print_options.PrintOptions
	kennung.Abbr
	maxKopf, maxSchwanz int
}

func (f *aligned) GetAligned() Aligned {
	return f
}

func (f *aligned) SetMaxKopfUndSchwanz(k, s int) {
	f.maxKopf = k
	f.maxSchwanz = s
}

func (f *aligned) WriteStringFormat(
	sw schnittstellen.WriterAndStringWriter,
	o *kennung.Kennung2,
) (n int64, err error) {
	var n1 int

	if f.Abbreviations.Hinweisen &&
		o.GetGattung() == gattung.Zettel &&
		!o.IsVirtual() {
		if err = f.AbbreviateKennung(o, o); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	// TODO move to kennung and avoid allocation
	h := kennung.Aligned(o, f.maxKopf, f.maxSchwanz)
	n1, err = sw.WriteString(h)
	n += int64(n1)

	return
}

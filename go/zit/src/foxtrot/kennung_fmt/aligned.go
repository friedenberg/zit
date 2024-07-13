package kennung_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

func MakeAligned(
	options erworben_cli_print_options.PrintOptions,
	ab ids.Abbr,
) *aligned {
	return &aligned{PrintOptions: options, Abbr: ab}
}

type aligned struct {
	erworben_cli_print_options.PrintOptions
	ids.Abbr
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
	sw interfaces.WriterAndStringWriter,
	o *ids.ObjectId,
) (n int64, err error) {
	var n1 int

	if f.Abbreviations.Hinweisen &&
		o.GetGenre() == genres.Zettel &&
		!o.IsVirtual() {
		if err = f.AbbreviateKennung(o, o); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	// TODO move to kennung and avoid allocation
	h := ids.Aligned(o, f.maxKopf, f.maxSchwanz)
	n1, err = sw.WriteString(h)
	n += int64(n1)

	return
}

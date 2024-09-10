package id_fmts

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	cli_print_options "code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

func MakeAligned(
	options cli_print_options.PrintOptions,
	ab ids.Abbr,
) *aligned {
	return &aligned{PrintOptions: options, Abbr: ab}
}

type aligned struct {
	cli_print_options.PrintOptions
	ids.Abbr
	maxHead, maxTail int
}

func (f *aligned) GetAligned() Aligned {
	return f
}

func (f *aligned) SetMaxKopfUndSchwanz(k, s int) {
	f.maxHead = k
	f.maxTail = s
}

func (f *aligned) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	o *ids.ObjectId,
) (n int64, err error) {
	rid := o.GetRepoId()
	var n2 int64

	switch o.GetGenre() {
	case genres.Zettel:
		if f.Abbreviations.Hinweisen && !o.IsVirtual() && rid == "" {
			if err = f.AbbreviateObjectId(o, o); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		fallthrough

	default:
		n2, err = f.writeId(sw, o)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f *aligned) writeId(
	sw interfaces.WriterAndStringWriter,
	o *ids.ObjectId,
) (n int64, err error) {
	var n1 int
	rid := o.GetRepoId()

	if len(rid) > 0 {
		n1, err = sw.WriteString(o.String())
		n += int64(n1)
	} else {
		// TODO move to object id and avoid allocation
		h := ids.Aligned(o, f.maxHead, f.maxTail)
		n1, err = sw.WriteString(h)
		n += int64(n1)
	}

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

package typ

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type FormatterValue struct {
	string
}

func (f FormatterValue) String() string {
	return f.string
}

func (f *FormatterValue) Set(v string) (err error) {
	v1 := strings.TrimSpace(strings.ToLower(v))
	switch v1 {
	case "action-names", "vim-syntax-type":
		f.string = v1

	default:
		err = objekte.MakeErrUnsupportedFormatterValue(v1, gattung.Typ)
		return
	}

	return
}

func (f *FormatterValue) FuncFormatter(
	out io.Writer,
	af schnittstellen.AkteIOFactory,
	agp schnittstellen.AkteGetterPutter[*Akte],
) schnittstellen.FuncIter[*sku.TransactedTyp] {
	switch f.string {
	case "action-names":
		f := MakeFormatterActionNames()

		return func(o *sku.TransactedTyp) (err error) {
			var akte *Akte

			if akte, err = agp.GetAkte(o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer agp.PutAkte(akte)

			if _, err = f.Format(out, akte); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "vim-syntax-type":
		f := MakeFormatterVimSyntaxType()

		return func(o *sku.TransactedTyp) (err error) {
			var akte *Akte

			if akte, err = agp.GetAkte(o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer agp.PutAkte(akte)

			if _, err = f.Format(out, akte); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	default:
		return func(_ *sku.TransactedTyp) (err error) {
			err = objekte.MakeErrUnsupportedFormatterValue(
				f.string,
				gattung.Typ,
			)
			return
		}
	}
}

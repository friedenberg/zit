package typ

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/charlie/gattung"
	"code.linenisgreat.com/zit-go/src/delta/typ_akte"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
	"code.linenisgreat.com/zit-go/src/juliett/objekte"
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
	agp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) schnittstellen.FuncIter[*sku.Transacted] {
	switch f.string {
	case "action-names":
		f := typ_akte.MakeFormatterActionNames()

		return func(o *sku.Transacted) (err error) {
			var akte *typ_akte.V0

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
		f := typ_akte.MakeFormatterVimSyntaxType()

		return func(o *sku.Transacted) (err error) {
			var akte *typ_akte.V0

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
		return func(_ *sku.Transacted) (err error) {
			err = objekte.MakeErrUnsupportedFormatterValue(
				f.string,
				gattung.Typ,
			)
			return
		}
	}
}

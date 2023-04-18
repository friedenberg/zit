package typ

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/golf/objekte"
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
	case "text", "objekte", "action-names", "vim-syntax-type":
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
) schnittstellen.FuncIter[*Transacted] {
	switch f.string {
	case "text":
		f := MakeFormatText(af)

		return func(o *Transacted) (err error) {
			if _, err = f.Format(out, &o.Objekte); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "action-names":
		f := MakeFormatterActionNames()

		return func(o *Transacted) (err error) {
			if _, err = f.Format(out, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "vim-syntax-type":
		f := MakeFormatterVimSyntaxType()

		return func(o *Transacted) (err error) {
			if _, err = f.Format(out, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	default:
		return func(_ *Transacted) (err error) {
			err = objekte.MakeErrUnsupportedFormatterValue(f.string, gattung.Typ)
			return
		}
	}
}

package erworben

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
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
	case "text", "objekte", "debug":
		f.string = v1

	default:
		err = errors.Errorf("unsupported format type: %s", v)
		return
	}

	return
}

func (f *FormatterValue) FuncFormatter(
	out io.Writer,
	af schnittstellen.AkteIOFactory,
) schnittstellen.FuncIter[*Transacted] {
	switch f.string {
	case "debug":
		return func(o *Transacted) (err error) {
			errors.Out().Printf("%#v", o)

			return
		}

	case "objekte":
		f := objekte.MakeFormat[Objekte, *Objekte]()

		return func(o *Transacted) (err error) {
			if _, err = f.Format(out, &o.Objekte); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "text":
		f := MakeFormatText(af)

		return func(o *Transacted) (err error) {
			if _, err = f.Format(out, &o.Objekte); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	default:
		return func(_ *Transacted) (err error) {
			return errors.Errorf("unsupported format for typen: %s", f.string)
		}
	}
}

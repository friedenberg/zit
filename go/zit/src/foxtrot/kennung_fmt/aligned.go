package kennung_fmt

import (
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

type Aligned struct {
	kennung.Abbr
	MaxKopf, MaxSchwanz int
	Padding             string
}

func (f *Aligned) SetMaxKopfUndSchwanz(k, s int) {
	f.MaxKopf = k
	f.MaxSchwanz = s
	f.Padding = strings.Repeat(" ", 5+k+s)
}

func (f *Aligned) WriteStringFormat(
	sw schnittstellen.WriterAndStringWriter,
	o *kennung.Kennung2,
) (n int64, err error) {
	var n1 int

	var k kennung.Kennung2

	if err = f.AbbreviateKennung(o, &k); err != nil {
		err = errors.Wrap(err)
		return
	}

	h := kennung.Aligned(&k, f.MaxKopf, f.MaxSchwanz)
	n1, err = sw.WriteString(h)
	n += int64(n1)

	return
}

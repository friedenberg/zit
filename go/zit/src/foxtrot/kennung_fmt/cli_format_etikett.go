package kennung_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type etikettenCliFormat struct{}

func MakeEtikettenCliFormat() (f *etikettenCliFormat) {
	f = &etikettenCliFormat{}

	return
}

func (f *etikettenCliFormat) WriteStringFormat(
	w schnittstellen.WriterAndStringWriter,
	k *kennung.Etikett,
) (n int64, err error) {
	var n1 int

	n1, err = w.WriteString(k.String())
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

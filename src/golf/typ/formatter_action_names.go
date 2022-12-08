package typ

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type FormatterActionNames struct {
	out io.Writer
}

func MakeFormatterActionNames(w io.Writer) FormatterActionNames {
	return FormatterActionNames{
		out: w,
	}
}

func (f FormatterActionNames) Encode(ct *Transacted) (n int64, err error) {
	for v, v1 := range ct.Objekte.Akte.Actions {
		var n1 int

		if n1, err = io.WriteString(f.out, fmt.Sprintf("%s\t%s\n", v, v1.Description)); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += int64(n1)
	}

	return
}

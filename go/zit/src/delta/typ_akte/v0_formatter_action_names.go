package typ_akte

import (
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type formatterActionNames struct{}

func MakeFormatterActionNames() *formatterActionNames {
	return &formatterActionNames{}
}

func (f formatterActionNames) Format(
	w io.Writer,
	ct *V0,
) (n int64, err error) {
	for v, v1 := range ct.Actions {
		var n1 int

		if n1, err = io.WriteString(
			w,
			fmt.Sprintf("%s\t%s\n", v, v1.Description),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += int64(n1)
	}

	return
}

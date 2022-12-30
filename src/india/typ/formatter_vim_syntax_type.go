package typ

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type formatterVimSyntaxType struct {
}

func MakeFormatterVimSyntaxType() *formatterVimSyntaxType {
	return &formatterVimSyntaxType{}
}

func (f formatterVimSyntaxType) Format(
	w io.Writer,
	ct *Transacted,
) (n int64, err error) {
	var n1 int

	if n1, err = io.WriteString(
		w,
		fmt.Sprintf("%s\n", ct.Objekte.Akte.VimSyntaxType),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	n += int64(n1)

	return
}

package typ_akte

import (
	"fmt"
	"io"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
)

type formatterVimSyntaxType struct{}

func MakeFormatterVimSyntaxType() *formatterVimSyntaxType {
	return &formatterVimSyntaxType{}
}

func (f formatterVimSyntaxType) Format(
	w io.Writer,
	ct *V0,
) (n int64, err error) {
	var n1 int

	if n1, err = fmt.Fprintln(
		w,
		ct.VimSyntaxType,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	n += int64(n1)

	return
}

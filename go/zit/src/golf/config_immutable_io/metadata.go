package config_immutable_io

import (
	"bufio"
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
)

type metadata struct {
	*ConfigLoaded
}

func (m metadata) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	if n, err = format.ReadLines(
		r,
		ohio.MakeLineReaderRepeat(
			ohio.MakeLineReaderKeyValues(
				map[string]interfaces.FuncSetString{
					"!": m.Type.Set,
				},
			),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (m metadata) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int
	n1, err = fmt.Fprintf(w, "! %s\n", m.Type.StringSansOp())
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

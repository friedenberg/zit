package ids

import (
	"bufio"
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
)

type TypeSetter interface {
	SetType(Type)
}

type TypedMetadataCoder[O interface {
	TypeGetter
	TypeSetter
}] struct{}

func (m TypedMetadataCoder[O]) DecodeFrom(
	object O,
	r1 io.Reader,
) (n int64, err error) {
	r := bufio.NewReader(r1)

	var t Type

	// TODO scan for type directly
	if n, err = format.ReadLines(
		r,
		ohio.MakeLineReaderRepeat(
			ohio.MakeLineReaderKeyValues(
				map[string]interfaces.FuncSetString{
					"!": t.Set,
				},
			),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	object.SetType(t)

	return
}

func (m TypedMetadataCoder[O]) EncodeTo(
	object O,
	w io.Writer,
) (n int64, err error) {
	var n1 int
	n1, err = fmt.Fprintf(w, "! %s\n", object.GetType().StringSansOp())
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

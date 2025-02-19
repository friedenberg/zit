package triple_hyphen_io

import (
	"bufio"
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
)

type TypedMetadataCoder[O any] struct{}

func (TypedMetadataCoder[O]) DecodeFrom(
	subject *TypedStruct[O],
	reader io.Reader,
) (n int64, err error) {
	bufferedReader := bufio.NewReader(reader)

	// TODO scan for type directly
	if n, err = format.ReadLines(
		bufferedReader,
		ohio.MakeLineReaderRepeat(
			ohio.MakeLineReaderKeyValues(
				map[string]interfaces.FuncSetString{
					"!": subject.Type.Set,
				},
			),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (TypedMetadataCoder[O]) EncodeTo(
	subject *TypedStruct[O],
	writer io.Writer,
) (n int64, err error) {
	var n1 int
	n1, err = fmt.Fprintf(writer, "! %s\n", subject.Type.StringSansOp())
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

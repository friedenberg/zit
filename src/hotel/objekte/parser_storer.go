package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
)

type ParserStorer[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
] struct {
	format gattung.Format[T, T1]
	owf    gattung.ObjekteWriterFactory
}

func MakeParserStorer[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
](
	owf gattung.ObjekteWriterFactory,
) *ParserStorer[T, T1] {
	return &ParserStorer[T, T1]{
		format: Format[T, T1]{},
		owf:    owf,
	}
}

func MakeParserStorerWithCustomFormat[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
](
	owf gattung.ObjekteWriterFactory,
	format gattung.Format[T, T1],
) *ParserStorer[T, T1] {
	if format == nil {
		panic("nil custom format for ParserStorer")
	}

	return &ParserStorer[T, T1]{
		format: format,
		owf:    owf,
	}
}

func (f ParserStorer[T, T1]) Parse(
	r io.Reader,
	o T1,
) (n int64, err error) {
	if n, err = f.format.Parse(r, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ow sha.WriteCloser

	if ow, err = f.owf.ObjekteWriter(o.GetGattung()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ow)

	errors.Todo(errors.P0, "must update sku with sha from writing to OW")
	if _, err = f.format.Format(ow, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

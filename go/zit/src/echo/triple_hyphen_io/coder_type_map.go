package triple_hyphen_io

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type TypedStruct[S any] struct {
	Type   *ids.Type
	Struct S
}

func (typedStruct *TypedStruct[S]) GetType() *ids.Type {
	if typedStruct.Type == nil {
		typedStruct.Type = &ids.Type{}
	}

	return typedStruct.Type
}

type CoderTypeMap[S any] map[string]interfaces.Coder[*TypedStruct[S]]

func (c CoderTypeMap[S]) DecodeFrom(
	subject *TypedStruct[S],
	reader io.Reader,
) (n int64, err error) {
	t := subject.GetType()
	coder, ok := c[t.String()]

	if !ok {
		err = errors.ErrorWithStackf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.DecodeFrom(subject, reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c CoderTypeMap[S]) EncodeTo(
	subject *TypedStruct[S],
	writer io.Writer,
) (n int64, err error) {
	t := subject.GetType()
	coder, ok := c[t.String()]

	if !ok {
		err = errors.ErrorWithStackf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.EncodeTo(subject, writer); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type DecoderTypeMapWithoutType[S any] map[string]interfaces.DecoderFrom[S]

func (c DecoderTypeMapWithoutType[S]) DecodeFrom(
	subject *TypedStruct[S],
	reader io.Reader,
) (n int64, err error) {
	t := subject.GetType()
	coder, ok := c[t.String()]

	if !ok {
		err = errors.ErrorWithStackf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.DecodeFrom(subject.Struct, reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type CoderTypeMapWithoutType[S any] map[string]interfaces.Coder[S]

func (c CoderTypeMapWithoutType[S]) DecodeFrom(
	subject *TypedStruct[S],
	reader io.Reader,
) (n int64, err error) {
	t := subject.GetType()
	coder, ok := c[t.String()]

	if !ok {
		err = errors.ErrorWithStackf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.DecodeFrom(subject.Struct, reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c CoderTypeMapWithoutType[S]) EncodeTo(
	subject *TypedStruct[S],
	writer io.Writer,
) (n int64, err error) {
	t := subject.Type
	coder, ok := c[t.String()]

	if !ok {
		err = errors.ErrorWithStackf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.EncodeTo(subject.Struct, writer); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

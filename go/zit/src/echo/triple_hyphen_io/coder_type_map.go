package triple_hyphen_io

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

// TODO rename and make clearer
type TypedStruct[O any] struct {
	Type   *ids.Type
	Struct O
}

func (typedStruct *TypedStruct[O]) GetType() *ids.Type {
	if typedStruct.Type == nil {
		typedStruct.Type = &ids.Type{}
	}

	return typedStruct.Type
}

type CoderTypeMap[O any] map[string]interfaces.Coder[*TypedStruct[O]]

func (c CoderTypeMap[O]) DecodeFrom(
	subject *TypedStruct[O],
	reader io.Reader,
) (n int64, err error) {
	t := subject.GetType()
	coder, ok := c[t.String()]

	if !ok {
		err = errors.Errorf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.DecodeFrom(subject, reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c CoderTypeMap[O]) EncodeTo(
	typedStruct *TypedStruct[O],
	writer io.Writer,
) (n int64, err error) {
	t := typedStruct.GetType()
	coder, ok := c[t.String()]

	if !ok {
		err = errors.Errorf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.EncodeTo(typedStruct, writer); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type TypedDecodersWithoutType[O any] map[string]interfaces.DecoderFrom[O]

func (c TypedDecodersWithoutType[O]) DecodeFrom(
	subject *TypedStruct[O],
	reader io.Reader,
) (n int64, err error) {
	t := subject.GetType()
	coder, ok := c[t.String()]

	if !ok {
		err = errors.Errorf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.DecodeFrom(subject.Struct, reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type TypedCodersWithoutType[O any] map[string]interfaces.Coder[O]

func (c TypedCodersWithoutType[O]) DecodeFrom(
	subject *TypedStruct[O],
	reader io.Reader,
) (n int64, err error) {
	t := subject.GetType()
	coder, ok := c[t.String()]

	if !ok {
		err = errors.Errorf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.DecodeFrom(subject.Struct, reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c TypedCodersWithoutType[O]) EncodeTo(
	subject *TypedStruct[O],
	writer io.Writer,
) (n int64, err error) {
	t := subject.Type
	coder, ok := c[t.String()]

	if !ok {
		err = errors.Errorf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.EncodeTo(subject.Struct, writer); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

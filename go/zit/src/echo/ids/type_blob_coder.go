package ids

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type TypeWithObject[O any] struct {
	Type   *Type
	Object O
}

func (typeWithObject TypeWithObject[O]) GetType() *Type {
	if typeWithObject.Type == nil {
		typeWithObject.Type = &Type{}
	}

	return typeWithObject.Type
}

type TypedCoders[O any] map[string]interfaces.Coder[*TypeWithObject[O]]

func (c TypedCoders[O]) DecodeFrom(
	object *TypeWithObject[O],
	reader io.Reader,
) (n int64, err error) {
	t := object.GetType()
	coder, ok := c[t.String()]

	if !ok {
		err = errors.Errorf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.DecodeFrom(object, reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c TypedCoders[O]) EncodeTo(
	object *TypeWithObject[O],
	writer io.Writer,
) (n int64, err error) {
	t := object.GetType()
	coder, ok := c[t.String()]

	if !ok {
		err = errors.Errorf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.EncodeTo(object, writer); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type TypedDecodersWithoutType[O any] map[string]interfaces.DecoderFrom[O]

func (c TypedDecodersWithoutType[O]) DecodeFrom(
	object *TypeWithObject[O],
	reader io.Reader,
) (n int64, err error) {
	t := object.GetType()
	coder, ok := c[t.String()]

	if !ok {
		err = errors.Errorf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.DecodeFrom(object.Object, reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type TypedCodersWithoutType[O any] map[string]interfaces.Coder[O]

func (c TypedCodersWithoutType[O]) DecodeFrom(
	object *TypeWithObject[O],
	reader io.Reader,
) (n int64, err error) {
	t := object.GetType()
	coder, ok := c[t.String()]

	if !ok {
		err = errors.Errorf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.DecodeFrom(object.Object, reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c TypedCodersWithoutType[O]) EncodeTo(
	subject *TypeWithObject[O],
	writer io.Writer,
) (n int64, err error) {
	t := subject.Type
	coder, ok := c[t.String()]

	if !ok {
		err = errors.Errorf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.EncodeTo(subject.Object, writer); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

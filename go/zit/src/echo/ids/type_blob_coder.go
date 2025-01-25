package ids

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type TypeGetter interface {
	GetType() Type
}

type TypedCoders[O TypeGetter] map[string]interfaces.Coder[O]

func (c TypedCoders[O]) DecodeFrom(object O, r io.Reader) (n int64, err error) {
	t := object.GetType()
	coder, ok := c[t.String()]

	if !ok {
		err = errors.Errorf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.DecodeFrom(object, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c TypedCoders[O]) EncodeTo(object O, w io.Writer) (n int64, err error) {
	t := object.GetType()
	coder, ok := c[t.String()]

	if !ok {
		err = errors.Errorf("no coders available for type: %q", t)
		return
	}

	if n, err = coder.EncodeTo(object, w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

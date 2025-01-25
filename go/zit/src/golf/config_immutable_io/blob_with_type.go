package config_immutable_io

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type BlobWithType[B any] struct {
	*ids.Type
	Coders map[string]interfaces.Coder[B]
}

func (c BlobWithType[B]) DecodeFrom(blob B, r io.Reader) (n int64, err error) {
	coder, ok := c.Coders[c.Type.String()]

	if !ok {
		err = errors.Errorf("no coders availabe for blob type: %q", c.Type)
		return
	}

	if n, err = coder.DecodeFrom(blob, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c BlobWithType[B]) EncodeTo(blob B, w io.Writer) (n int64, err error) {
	coder, ok := c.Coders[c.Type.String()]

	if !ok {
		err = errors.Errorf("no coders availabe for blob type: %q", c.Type)
		return
	}

	if n, err = coder.EncodeTo(blob, w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

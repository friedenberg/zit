package ohio

import (
	"bytes"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func WriteKeySpaceValueNewlineString(
	w io.Writer,
	key, value string,
) (n int, err error) {
	return WriteStrings(w, key, " ", value, "\n")
}

func WriteKeySpaceValueNewline(
	w io.Writer,
	key string, value []byte,
) (n int64, err error) {

	var (
		n1 int64
		b  *bytes.Buffer
	)

	b = bytes.NewBufferString(key)
	n1, err = b.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	b = bytes.NewBufferString(" ")

	n1, err = b.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	b = bytes.NewBuffer(value)

	n1, err = b.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	b = bytes.NewBufferString("\n")

	n1, err = b.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func WriteStrings(
	w io.Writer,
	ss ...string,
) (n int, err error) {
	for _, s := range ss {
		var n1 int

		n1, err = io.WriteString(w, s)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

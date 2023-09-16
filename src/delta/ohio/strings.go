package ohio

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func WriteKeySpaceValueNewline(
	w io.Writer,
	key, value string,
) (n int, err error) {
	return WriteStrings(w, key, " ", value, "\n")
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

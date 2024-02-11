package ohio

import (
	"bytes"
	"io"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
)

func WriteKeySpaceValueNewlineString(
	w io.Writer,
	key, value string,
) (n int, err error) {
	return WriteStrings(w, key, " ", value, "\n")
}

func WriteKeySpaceValueNewlineWritersTo(
	w io.Writer,
	key, value io.WriterTo,
) (n int64, err error) {
	var (
		n1 int64
		n2 int
	)

	n1, err = key.WriteTo(w)
	n += n1

	if err != nil {
		return
	}

	n2, err = w.Write([]byte{' '})
	n += int64(n2)

	if err != nil {
		return
	}

	n1, err = value.WriteTo(w)
	n += n1

	if err != nil {
		return
	}

	n2, err = w.Write([]byte{'\n'})
	n += int64(n2)

	if err != nil {
		return
	}

	return
}

func WriteKeySpaceValueNewline(
	w io.Writer,
	key string, value []byte,
) (n int64, err error) {
	var (
		n1 int64
		sr *strings.Reader
		br *bytes.Reader
	)

	sr = strings.NewReader(key)
	n1, err = sr.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	sr = strings.NewReader(" ")

	n1, err = sr.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	br = bytes.NewReader(value)

	n1, err = br.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	sr = strings.NewReader("\n")

	n1, err = sr.WriteTo(w)
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

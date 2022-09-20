package gattung

import (
	"bufio"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func ConfirmTypeFromReader(t Gattung, r *bufio.Reader) (err error) {
	var t1 Gattung

	if t1, err = FromReader(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	if t1 != t {
		err = errors.Normal(
			ErrWrongType{
				ExpectedType: t,
				ActualType:   t,
			},
		)
	}

	return
}

func FromReader(r *bufio.Reader) (t Gattung, err error) {
	var line string

	if line, err = r.ReadString('\n'); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = t.Set(line); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

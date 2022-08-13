package zk_types

import (
	"bufio"

	"github.com/friedenberg/zit/bravo/errors"
)

func ConfirmTypeFromReader(t Type, r *bufio.Reader) (err error) {
	var t1 Type

	if t1, err = FromReader(r); err != nil {
		err = errors.Error(err)
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

func FromReader(r *bufio.Reader) (t Type, err error) {
	var line string

	if line, err = r.ReadString('\n'); err != nil {
		err = errors.Error(err)
		return
	}

	if err = t.Set(line); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

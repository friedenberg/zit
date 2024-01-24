package sha

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
)

var ErrIsNull = errors.New("sha is null")

func MakeErrIsNull(s *Sha) error {
	if s.IsNull() {
		return errors.WrapN(1, ErrIsNull)
	}

	return nil
}

type errLength [2]int

func makeErrLength(expected, actual int) error {
	if expected != actual {
		return errLength{expected, actual}
	} else {
		return nil
	}
}

func (e errLength) Error() string {
	return fmt.Sprintf("expected %d but got %d", e[0], e[1])
}

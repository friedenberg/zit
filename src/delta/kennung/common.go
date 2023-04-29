package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func LeftSubtract[T schnittstellen.Stringer, TPtr schnittstellen.StringSetterPtr[T]](
	a, b T,
) (c T, err error) {
	if err = TPtr(&c).Set(strings.TrimPrefix(a.String(), b.String())); err != nil {
		err = errors.Wrapf(err, "'%s' - '%s'", a, b)
		return
	}

	return
}

func Contains[T schnittstellen.Stringer](a, b T) bool {
	if len(b.String()) > len(a.String()) {
		return false
	}

	return strings.HasPrefix(a.String(), b.String())
}

func Less[T schnittstellen.Stringer](a, b T) bool {
	return a.String() < b.String()
}

func IsEmpty[T schnittstellen.Stringer](a T) bool {
	return len(a.String()) == 0
}

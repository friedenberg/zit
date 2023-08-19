package store_objekten

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type ErrAkteExists struct {
	Akte sha.Sha
	zettel.MutableSet
}

func (e ErrAkteExists) Is(target error) bool {
	_, ok := target.(ErrAkteExists)
	return ok
}

func (e ErrAkteExists) Error() string {
	return fmt.Sprintf(
		"zettelen already exist with akte:\n%s\n%v",
		e.Akte,
		e.MutableSet,
	)
}

type ErrExternalAkteExtensionMismatch struct {
	Expected string
	Actual   kennung.FD
}

func (e ErrExternalAkteExtensionMismatch) Is(target error) bool {
	_, ok := target.(ErrExternalAkteExtensionMismatch)
	return ok
}

func (e ErrExternalAkteExtensionMismatch) Error() string {
	return fmt.Sprintf(
		"expected extension %q but got %q",
		e.Expected,
		e.Actual,
	)
}

func IsErrExternalAkteExtensionMismatch(err error) bool {
	return errors.Is(err, ErrExternalAkteExtensionMismatch{})
}

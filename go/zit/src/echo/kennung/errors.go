package kennung

import (
	"errors"
	"fmt"
)

type ErrInvalidKennung string

func (e ErrInvalidKennung) Error() string {
	return fmt.Sprintf("invalid kennung: %q", string(e))
}

func (e ErrInvalidKennung) Is(err error) (ok bool) {
	_, ok = err.(ErrInvalidKennung)
	return
}

func IsErrInvalid(err error) bool {
	return errors.Is(err, ErrInvalidKennung(""))
}

type errInvalidSigil string

func (e errInvalidSigil) Error() string {
	return fmt.Sprintf("invalid sigil: %q", string(e))
}

func (e errInvalidSigil) Is(err error) (ok bool) {
	_, ok = err.(errInvalidSigil)
	return
}

func IsErrInvalidSigil(err error) bool {
	return errors.Is(err, errInvalidSigil(""))
}

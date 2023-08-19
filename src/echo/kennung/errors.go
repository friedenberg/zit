package kennung

import (
	"errors"
	"fmt"
)

type errInvalidKennung string

func (e errInvalidKennung) Error() string {
	return fmt.Sprintf("invalid kennung: %q", string(e))
}

func (e errInvalidKennung) Is(err error) (ok bool) {
	_, ok = err.(errInvalidKennung)
	return
}

func IsErrInvalid(err error) bool {
	return errors.Is(err, errInvalidKennung(""))
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

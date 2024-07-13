package ids

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type ErrInvalidId string

func (e ErrInvalidId) Error() string {
	return fmt.Sprintf("invalid kennung: %q", string(e))
}

func (e ErrInvalidId) Is(err error) (ok bool) {
	_, ok = err.(ErrInvalidId)
	return
}

func IsErrInvalid(err error) bool {
	return errors.Is(err, ErrInvalidId(""))
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

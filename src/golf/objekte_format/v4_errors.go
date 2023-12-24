package objekte_format

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
)

var (
	errV4ExpectedNewline           = errors.New("expected newline")
	ErrV4ExpectedSpaceSeparatedKey = errors.New("expected space separated key")
	errV4EmptyKey                  = errors.New("empty key")
	errV4KeysNotSorted             = errors.New("keys not sorted")
	errV4InvalidKey                = errors.New("invalid key")
)

func makeErrWithBytes(err error, bs []byte) error {
	return errors.WrapN(1, errWithBytes{error: err, bytes: bs})
}

type errWithBytes struct {
	error
	bytes []byte
}

func (ewb errWithBytes) Error() string {
	return fmt.Sprintf("%s: %s", ewb.error, ewb.bytes)
}

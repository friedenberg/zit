package objekte_format

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

var (
	errV4ExpectedNewline           = errors.New("expected newline")
	ErrV4ExpectedSpaceSeparatedKey = errors.New("expected space separated key")
	errV4EmptyKey                  = errors.New("empty key")
	errV4KeysNotSorted             = errors.New("keys not sorted")
	errV4InvalidKey                = errors.New("invalid key")
)

func makeErrWithBytes(err error, bs []byte) error {
	if ui.IsVerbose() {
		return errors.WrapN(1, errWithBytes{error: err, bytes: bs})
	}

	return err
}

type errWithBytes struct {
	error
	bytes []byte
}

func (ewb errWithBytes) Error() string {
	return fmt.Sprintf("%s: %s", ewb.error, ewb.bytes)
}

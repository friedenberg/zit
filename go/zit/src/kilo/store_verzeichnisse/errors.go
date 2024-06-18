package store_verzeichnisse

import "code.linenisgreat.com/zit/go/zit/src/alfa/errors"

var errConcurrentPageAccess = errors.New("concurrent page access")

func MakeErrConcurrentPageAccess() error {
	return errors.WrapN(2, errConcurrentPageAccess)
}

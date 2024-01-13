package store_verzeichnisse

import "github.com/friedenberg/zit/src/alfa/errors"

var errConcurrentPageAccess = errors.New("concurrent page access")

func MakeErrConcurrentPageAccess() error {
	return errors.WrapN(2, errConcurrentPageAccess)
}

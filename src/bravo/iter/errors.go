package iter

import "github.com/friedenberg/zit/src/alfa/errors"

var errStopIteration = errors.New("stop iteration")

func MakeErrStopIteration() error {
	if errors.IsVerbose() {
		return errors.WrapN(2, errStopIteration)
	} else {
		return errStopIteration
	}
}

func IsStopIteration(err error) bool {
	if errors.Is(err, errStopIteration) {
		errors.Log().Printf("stopped iteration at %s", err)
		return true
	}

	return false
}

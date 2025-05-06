package checkout_mode

import "code.linenisgreat.com/zit/go/zit/src/alfa/errors"

type errInvalidCheckoutMode error

func MakeErrInvalidCheckoutModeMode(mode Mode) errInvalidCheckoutMode {
	return errors.WrapN(
		1,
		errInvalidCheckoutMode(
			errors.ErrorWithStackf("invalid checkout mode: %s", mode),
		),
	)
}

func MakeErrInvalidCheckoutMode(err error) errInvalidCheckoutMode {
	return errInvalidCheckoutMode(err)
}

func IsErrInvalidCheckoutMode(err error) bool {
	return errors.Is(err, errInvalidCheckoutMode(nil))
}

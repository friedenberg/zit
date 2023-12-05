package errors

var (
	ErrStopIteration = New("stop iteration")
	ErrFalse         = New("false")
	ErrTrue          = New("true")
)

func MakeErrStopIteration() error {
	if IsVerbose() {
		return WrapN(2, ErrStopIteration)
	} else {
		return ErrStopIteration
	}
}

func IsStopIteration(err error) bool {
	if Is(err, ErrStopIteration) {
		// errors.Log().Printf("stopped iteration at %s", err)
		return true
	}

	return false
}

func IsErrFalse(err error) bool {
	if Is(err, ErrFalse) {
		return true
	}

	return false
}

func IsErrTrue(err error) bool {
	if Is(err, ErrTrue) {
		return true
	}

	return false
}

package errors

var (
	ErrStopIteration = New("stop iteration")
	ErrFalse         = New("false")
	ErrTrue          = New("true")
)

func MakeErrStopIteration() error {
	return ErrStopIteration
}

func IsStopIteration(err error) bool {
	if Is(err, ErrStopIteration) {
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

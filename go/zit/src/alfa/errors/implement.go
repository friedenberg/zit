package errors

var errImplement = New("not implemented")

func Implement() (err error) {
	// Err().Caller(1, "%s", errImplement)
	return WrapN(1, errImplement)
}

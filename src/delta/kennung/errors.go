package kennung

import "errors"

var (
	ErrInvalid = errors.New("invalid kennung")
)

func IsErrInvalid(err error) bool {
	return errors.Is(err, ErrInvalid)
}

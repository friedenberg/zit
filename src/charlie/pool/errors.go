package pool

import "errors"

var ErrDoNotRepool = errors.New("do not repool")

func IsDoNotRepool(err error) bool {
	return errors.Is(err, ErrDoNotRepool)
}

package pool

import "code.linenisgreat.com/zit/go/zit/src/alfa/errors"

var ErrDoNotRepool = errors.New("do not repool")

func IsDoNotRepool(err error) bool {
	return errors.Is(err, ErrDoNotRepool)
}

package user_ops

import (
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
)

type OpenAkten struct {
}

func (c OpenAkten) RunMany(s store_with_lock.Store, zs zettel_transacted.Set) (err error) {
	// if len(zs.SetNamed) == 0 {
	// 	return
	// }

	// if err = open_file_guard.Open(args...); err != nil {
	// 	err = errors.Errorf("%q: %s", args, err)
	// 	return
	// }

	return
}

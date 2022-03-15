package commands

import "github.com/friedenberg/zit/india/store_with_lock"

type WithShas interface {
	RunWithShas(store store_with_lock.Store, shas ..._Sha) error
}

type withShas struct {
	WithShas
}

func (c withShas) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	shas := make([]_Sha, len(args))

	for i, arg := range args {
		var sha _Sha

		if err = sha.Set(arg); err != nil {
			err = _Error(err)
			return
		}

		shas[i] = sha
	}

	if err = c.RunWithShas(store, shas...); err != nil {
		err = _Error(err)
		return
	}

	return
}

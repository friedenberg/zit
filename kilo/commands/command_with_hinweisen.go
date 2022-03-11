package commands

import "github.com/friedenberg/zit/india/store_with_lock"

type CommandWithHinweisen interface {
	RunWithHinweisen(_Umwelt, _Zettels, ..._Hinweis) error
}

type commandWithHinweisen struct {
	CommandWithHinweisen
}

func (c commandWithHinweisen) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	ids := make([]_Hinweis, len(args))

	for i, arg := range args {
		var h _Hinweis

		if h, err = _MakeBlindHinweis(arg); err != nil {
			err = _Error(err)
			return
		}

		ids[i] = h
	}

	c.RunWithHinweisen(store.Umwelt, store.Zettels(), ids...)

	return
}

package commands

import "github.com/friedenberg/zit/india/store_with_lock"

type CommandWithId interface {
	RunWithId(_Umwelt, _Zettels, ..._Id) error
}

type commandWithId struct {
	CommandWithId
}

func (c commandWithId) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	ids := make([]_Id, len(args))

	for i, arg := range args {
		var id _Id
		var sha _Sha

		if err = sha.Set(arg); err != nil {
			if id, err = _MakeBlindHinweis(arg); err != nil {
				err = _Error(err)
				return
			}
		} else {
			id = sha
		}

		ids[i] = id
	}

	c.RunWithId(store.Umwelt, store.Zettels(), ids...)

	return
}

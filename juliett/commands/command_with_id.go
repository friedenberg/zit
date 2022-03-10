package commands

type CommandWithId interface {
	RunWithId(_Umwelt, _Zettels, ..._Id) error
}

type commandWithId struct {
	CommandWithId
}

func (c commandWithId) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
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

	c.RunWithId(u, zs, ids...)

	return
}

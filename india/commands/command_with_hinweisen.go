package commands

type CommandWithHinweisen interface {
	RunWithHinweisen(_Umwelt, _Zettels, ..._Hinweis) error
}

type commandWithHinweisen struct {
	CommandWithHinweisen
}

func (c commandWithHinweisen) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	ids := make([]_Hinweis, len(args))

	for i, arg := range args {
		var h _Hinweis

		if h, err = _MakeBlindHinweis(arg); err != nil {
			err = _Error(err)
			return
		}

		ids[i] = h
	}

	c.RunWithHinweisen(u, zs, ids...)

	return
}

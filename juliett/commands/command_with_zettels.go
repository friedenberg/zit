package commands

type CommandWithZettels interface {
	RunWithZettels(_Umwelt, _Zettels, ...string) error
}

type commandWithZettels struct {
	CommandWithZettels
}

func (c commandWithZettels) Run(u _Umwelt, args ...string) (err error) {
	var age _Age

	if age, err = u.Age(); err != nil {
		err = _Error(err)
		return
	}

	var zs _Zettels

	if zs, err = _NewZettels(u, age); err != nil {
		err = _Error(err)
		return
	}

	defer _PanicIfError(zs.Flush)

	if err = c.RunWithZettels(u, zs, args...); err != nil {
		err = _Error(err)
		return
	}

	return
}

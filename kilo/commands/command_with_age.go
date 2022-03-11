package commands

type CommandWithAge interface {
	RunWithAge(_Umwelt, _Age, ...string) error
}

type commandWithAge struct {
	CommandWithAge
}

func (c commandWithAge) Run(u _Umwelt, args ...string) (err error) {
	var age _Age

	if age, err = u.Age(); err != nil {
		err = _Error(err)
		return
	}

	if err = c.RunWithAge(u, age, args...); err != nil {
		err = _Error(err)
		return
	}

	return
}

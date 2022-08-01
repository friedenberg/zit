package user_ops

type OpenFiles struct {
}

func (c OpenFiles) Run(args ...string) (err error) {
	if len(args) == 0 {
		return
	}

	if err = _OpenFiles(args...); err != nil {
		err = _Errorf("%q: %s", args, err)
		return
	}

	return
}

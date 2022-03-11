package user_ops

type Checkout struct {
	Options _ZettelsCheckinOptions
	Umwelt  _Umwelt
	Store   _Store
}

type CheckoutResults struct {
	Zettelen      []_ZettelCheckedOut
	FilesZettelen []string
	FilesAkten    []string
}

func (c Checkout) Run(args ...string) (results CheckoutResults, err error) {
	if results.Zettelen, err = c.Store.Checkout(c.Options, args...); err != nil {
		err = _Error(err)
		return
	}

	results.FilesZettelen = make([]string, 0, len(results.Zettelen))
	results.FilesAkten = make([]string, 0)

	for _, z := range results.Zettelen {
		results.FilesZettelen = append(results.FilesZettelen, z.External.Path)

		if z.External.AktePath != "" {
			results.FilesAkten = append(results.FilesAkten, z.External.AktePath)
		}
	}

	return
}

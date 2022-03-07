package commands

import (
	"bufio"
	"flag"
)

type CatAlfred struct {
	Type _Type
	Command
}

func init() {
	registerCommand(
		"cat-alfred",
		func(f *flag.FlagSet) Command {
			c := &CatAlfred{
				Type: _TypeUnknown,
			}

			c.Command = commandWithZettels{c}

			f.Var(&c.Type, "type", "ObjekteType")

			return c
		},
	)
}

func (c CatAlfred) HandleError(u _Umwelt, in error) {
	wo := bufio.NewWriter(u.Out)
	defer wo.Flush()

	var aw _AlfredWriter

	var err error

	if aw, err = _AlfredNewWriter(u.Out); err != nil {
		_PanicIfError(err)
		return
	}

	aw.WriteError(in)
	_PanicIfError(aw.Close())
}

func (c CatAlfred) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	//this command does its own error handling
	defer func() {
		err = nil
	}()

	wo := bufio.NewWriter(u.Out)
	defer wo.Flush()

	var aw _AlfredWriter

	if aw, err = _AlfredNewWriter(u.Out); err != nil {
		err = _Error(err)
		return
	}

	defer _PanicIfError(aw.Close)

	defer func() {
		aw.WriteError(err)
	}()

	switch c.Type {
	case _TypeEtikett:
		var ea []_Etikett

		if ea, err = zs.Etiketten().All(); err != nil {
			err = _Error(err)
			return
		}

		for _, e := range ea {
			aw.WriteEtikett(e)
		}

	case _TypeZettel:

		var all map[string]_NamedZettel

		if all, err = zs.All(); err != nil {
			err = _Error(err)
			return
		}

		for _, z := range all {
			aw.WriteZettel(z)
		}

	case _TypeAkte:

	case _TypeHinweis:

		var all map[string]_NamedZettel

		if all, err = zs.All(); err != nil {
			err = _Error(err)
			return
		}

		for _, z := range all {
			aw.WriteZettel(z)
		}

	default:
		err = _Errorf("unsupported objekte type: %s", c.Type)
		return
	}

	return
}

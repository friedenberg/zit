package commands

import (
	"bufio"
	"flag"
)

type CatAlfred struct {
	Type _Type
}

func init() {
	registerCommand(
		"cat-alfred",
		func(f *flag.FlagSet) Command {
			c := &CatAlfred{
				Type: _TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithZettels{c}
		},
	)
}

func (c CatAlfred) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	wo := bufio.NewWriter(u.Out)
	defer wo.Flush()

	var wa _AlfredZettelsWriter

	if wa, err = _AlfredZettelsNewWriter(u.Out); err != nil {
		err = _Error(err)
		return
	}

	defer _PanicIfError(wa.Close)

	switch c.Type {
	case _TypeEtikett:
		var ea []_Etikett

		if ea, err = zs.Etiketten().All(); err != nil {
			err = _Error(err)
			return
		}

		for _, e := range ea {
			wa.WriteEtikett(e)
		}

	case _TypeZettel:

		var all map[string]_NamedZettel

		if all, err = zs.All(); err != nil {
			err = _Error(err)
			return
		}

		for _, z := range all {
			wa.WriteZettel(z)
		}

	case _TypeAkte:

	case _TypeHinweis:

		var all map[string]_NamedZettel

		if all, err = zs.All(); err != nil {
			err = _Error(err)
			return
		}

		for _, z := range all {
			wa.WriteZettel(z)
		}

	default:
		err = _Errorf("unsupported objekte type: %s", c.Type)
		return
	}

	return
}

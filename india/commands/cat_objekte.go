package commands

import (
	"flag"
)

type CatObjekte struct {
	Type _Type
}

func init() {
	registerCommand(
		"cat-objekte",
		func(f *flag.FlagSet) Command {
			c := &CatObjekte{
				Type: _TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithZettels{c}
		},
	)
}

func (c CatObjekte) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	switch c.Type {

	case _TypeAkte:
		return c.akten(u, zs, args...)

	case _TypeZettel:
		return c.zettelen(u, zs, args...)

	default:
		err = _Errorf("unsupported objekte type: %s", c.Type)
		return
	}

	return
}

func (c CatObjekte) akten(u _Umwelt, zs _Zettels, args ...string) (err error) {
	for _, arg := range args {
		var sb _Sha

		if err = sb.Set(arg); err != nil {
			err = _Error(err)
			return
		}

		p := u.DirZit("Objekte", "Akte")

		if sb, err = sb.Glob(p); err != nil {
			err = _Error(err)
			return
		}

		if err = _ObjekteRead(u.Out, zs.Age(), _IdPath(sb, p)); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}

func (c CatObjekte) zettelen(u _Umwelt, zs _Zettels, args ...string) (err error) {
	for _, arg := range args {
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

		var z _NamedZettel

		if z, err = zs.Read(id); err != nil {
			err = _Error(err)
			return
		}

		f := _StoredZettelFormatsObjekte{}

		if _, err = f.WriteTo(z.Stored, u.Out); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}

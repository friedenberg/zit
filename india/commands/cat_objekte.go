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

			return commandWithZettels{commandWithId{c}}
		},
	)
}

func (c CatObjekte) RunWithId(u _Umwelt, zs _Zettels, ids ..._Id) (err error) {
	switch c.Type {

	case _TypeAkte:
		return c.akten(u, zs, ids...)

	case _TypeZettel:
		return c.zettelen(u, zs, ids...)

	default:
		err = _Errorf("unsupported objekte type: %s", c.Type)
		return
	}

	return
}

func (c CatObjekte) akten(u _Umwelt, zs _Zettels, ids ..._Id) (err error) {
	for _, id := range ids {
		var sb _Sha

		switch i := id.(type) {
		case _Sha:
			sb = i

		case _Hinweis:
			var named _NamedZettel

			if named, err = zs.Read(i); err != nil {
				err = _Error(err)
				return
			}

			sb = named.Zettel.Akte

		default:
			err = _Errorf("unsupported id type: %q", i)
			return
		}

		p := u.DirAkte()

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

func (c CatObjekte) zettelen(u _Umwelt, zs _Zettels, ids ..._Id) (err error) {
	for _, id := range ids {
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

package commands

import (
	"flag"
)

type Revert struct {
	Type _Type
}

func init() {
	registerCommand(
		"revert",
		func(f *flag.FlagSet) Command {
			c := &Revert{
				Type: _TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithZettels{c}
		},
	)
}

func (c Revert) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	u.Lock.Lock()
	defer _PanicIfError(u.Lock.Unlock())

	switch c.Type {
	case _TypeZettel:
		hins := make([]_Hinweis, len(args))

		for i, arg := range args {
			if hins[i], err = _MakeBlindHinweis(arg); err != nil {
				err = _Error(err)
				return
			}

			if _, err = zs.Revert(hins[i]); err != nil {
				err = _Error(err)
				return
			}
		}

	default:
		_Errf("objekte type %s does not support reverts currently\n", c.Type)
	}

	return
}

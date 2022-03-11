package commands

import (
	"flag"
	"os"
)

type Status struct {
}

func init() {
	registerCommand(
		"status",
		func(f *flag.FlagSet) Command {
			c := &Status{}

			return commandWithZettels{c}
		},
	)
}

func (c Status) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	if len(args) > 0 {
		_Errf("args provided will be ignored")
	}

	var cwd string

	if cwd, err = os.Getwd(); err != nil {
		err = _Error(err)
		return
	}

	var hins []string

	if hins, err = zs.GetPossibleZettels(cwd); err != nil {
		err = _Error(err)
		return
	}

	var daZees map[_Hinweis]_ExternalZettel

	options := _ZettelsCheckinOptions{
		IncludeAkte: true,
		Format:      _ZettelFormatsText{},
	}

	if daZees, err = zs.ReadExternal(options, hins...); err != nil {
		err = _Error(err)
		return
	}

	for h, z := range daZees {
		var named _NamedZettel

		if named, err = zs.Read(h); err != nil {
			err = _Error(err)
			return
		}

		if named.Zettel.Equals(z.Zettel) {
			continue
		}

		_Outf("[%s] (different)\n", h)
	}

	return
}

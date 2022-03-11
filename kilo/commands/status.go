package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/india/store_with_lock"
)

type Status struct {
}

func init() {
	registerCommand(
		"status",
		func(f *flag.FlagSet) Command {
			c := &Status{}

			return commandWithLockedStore{c}
		},
	)
}

func (c Status) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	if len(args) > 0 {
		_Errf("args provided will be ignored")
	}

	var cwd string

	if cwd, err = os.Getwd(); err != nil {
		err = _Error(err)
		return
	}

	var hins []string

	if hins, err = store.Zettels().GetPossibleZettels(cwd); err != nil {
		err = _Error(err)
		return
	}

	var daZees map[_Hinweis]_ExternalZettel

	options := _ZettelsCheckinOptions{
		IncludeAkte: true,
		Format:      _ZettelFormatsText{},
	}

	if daZees, err = store.Zettels().ReadExternal(options, hins...); err != nil {
		err = _Error(err)
		return
	}

	for h, z := range daZees {
		var named _NamedZettel

		if named, err = store.Zettels().Read(h); err != nil {
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

package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/india/store_with_lock"
)

type Clean struct {
}

func init() {
	registerCommand(
		"clean",
		func(f *flag.FlagSet) Command {
			c := &Clean{}

			return commandWithLockedStore{c}
		},
	)
}

func (c Clean) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	if len(args) > 0 {
		_Errf("args provided will be ignored")
	}

	var cwd string

	if cwd, err = os.Getwd(); err != nil {
		err = _Error(err)
		return
	}

	var hins []string

	//TODO move to user_ops
	if hins, err = store.Zettels().GetPossibleZettels(cwd); err != nil {
		err = _Error(err)
		return
	}

	var daZees map[_Hinweis]_ExternalZettel

	options := _ZettelsCheckinOptions{
		IncludeAkte: true,
		Format:      _ZettelFormatsText{},
	}

	//TODO move to user_ops
	if daZees, err = store.Zettels().ReadExternal(options, hins...); err != nil {
		err = _Error(err)
		return
	}

	toDelete := make([]_ExternalZettel, 0, len(daZees))
	filesToDelete := make([]string, 0, len(daZees))

	for h, z := range daZees {
		var named _NamedZettel

		if named, err = store.Zettels().Read(h); err != nil {
			err = _Error(err)
			return
		}

		if !named.Zettel.Equals(z.Zettel) {
			continue
		}

		toDelete = append(toDelete, z)
		filesToDelete = append(filesToDelete, z.Path)

		if z.AktePath != "" {
			filesToDelete = append(filesToDelete, z.AktePath)
		}
	}

	if store.Konfig.DryRun {
		for _, z := range toDelete {
			_Outf("[%s] (would delete)\n", z.Hinweis)
		}

		return
	}

	if err = _DeleteFilesAndDirs(filesToDelete...); err != nil {
		err = _Error(err)
		return
	}

	for _, z := range toDelete {
		_Outf("[%s] (checkout deleted)\n", z.Hinweis)
	}

	return
}

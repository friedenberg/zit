package commands

import (
	"flag"
	"os"
)

type Clean struct {
}

func init() {
	registerCommand(
		"clean",
		func(f *flag.FlagSet) Command {
			c := &Clean{}

			return commandWithZettels{c}
		},
	)
}

func (c Clean) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
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

	toDelete := make([]_ExternalZettel, 0, len(daZees))
	filesToDelete := make([]string, 0, len(daZees))

	for h, z := range daZees {
		var named _NamedZettel

		if named, err = zs.Read(h); err != nil {
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

	if u.Konfig.DryRun {
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

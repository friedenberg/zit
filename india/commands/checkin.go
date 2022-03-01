package commands

import (
	"flag"
	"os"
)

type Checkin struct {
	Delete     bool
	IgnoreAkte bool
	All        bool
}

func init() {
	registerCommand(
		"checkin",
		func(f *flag.FlagSet) Command {
			c := &Checkin{}

			f.BoolVar(&c.Delete, "delete", false, "the checked-out file")
			f.BoolVar(&c.IgnoreAkte, "ignore-akte", false, "do not change the akte")
			f.BoolVar(&c.All, "all", false, "")

			return commandWithZettels{c}
		},
	)
}

func (c Checkin) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	u.Lock.Lock()
	defer _PanicIfError(u.Lock.Unlock())

	if c.All {
		if len(args) > 0 {
			_Errf("Ignoring args because -all is set\n")
		}

		var cwd string

		if cwd, err = os.Getwd(); err != nil {
			err = _Error(err)
			return
		}

		if args, err = zs.GetPossibleZettels(cwd); err != nil {
			err = _Error(err)
			return
		}
	}

	options := _ZettelsCheckinOptions{
		IncludeAkte: !c.IgnoreAkte,
		Format:      _ZettelFormatsText{},
	}

	var daZees map[_Hinweis]_ZettelCheckedOut

	if daZees, err = zs.Checkin(options, args...); err != nil {
		err = _Error(err)
		return
	}

	if c.Delete {
		if err = c.deleteCheckouts(zs, daZees); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}

//TODO combine with clean command in zettel store
func (c Checkin) deleteCheckouts(zs _Zettels, daZees map[_Hinweis]_ZettelCheckedOut) (err error) {
	toDelete := make([]_ExternalZettel, 0, len(daZees))
	filesToDelete := make([]string, 0, len(daZees))

	for _, z := range daZees {
		if !z.Internal.Zettel.Equals(z.External.Zettel) {
			continue
		}

		toDelete = append(toDelete, z.External)
		filesToDelete = append(filesToDelete, z.External.Path)

		if z.External.AktePath != "" {
			filesToDelete = append(filesToDelete, z.External.AktePath)
		}
	}

	if err = _DeleteFilesAndDirs(filesToDelete...); err != nil {
		err = _Error(err)
		return
	}

	for _, z := range toDelete {
		_Outf("[%s] (checkout deleted)\n", z.Hinweis)
	}

	return

	return
}

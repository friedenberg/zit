package commands

import (
	"flag"

	"github.com/friedenberg/zit/juliett/user_ops"
)

type Clean struct {
}

func init() {
	registerCommand(
		"clean",
		func(f *flag.FlagSet) Command {
			c := &Clean{}

			return c
		},
	)
}

func (c Clean) Run(u _Umwelt, args ...string) (err error) {
	if len(args) > 0 {
		_Errf("args provided will be ignored")
	}

	getPossibleOp := user_ops.GetPossibleZettels{
		Umwelt: u,
	}

	var getPossibleResults user_ops.GetPossibleZettelsResults

	if getPossibleResults, err = getPossibleOp.Run(); err != nil {
		err = _Error(err)
		return
	}

	args = getPossibleResults.Hinweisen

	checkinOptions := _ZettelsCheckinOptions{
		IncludeAkte: true,
		Format:      _ZettelFormatsText{},
	}

	var readResults user_ops.ReadCheckedOutResults

	readOp := user_ops.ReadCheckedOut{
		Umwelt:  u,
		Options: checkinOptions,
	}

	if readResults, err = readOp.Run(args...); err != nil {
		err = _Error(err)
		return
	}

	toDelete := make([]_ExternalZettel, 0, len(readResults.Zettelen))
	filesToDelete := make([]string, 0, len(readResults.Zettelen))

	for _, z := range readResults.Zettelen {
		if !z.Internal.Zettel.Equals(z.External.Zettel) {
			continue
		}

		toDelete = append(toDelete, z.External)
		filesToDelete = append(filesToDelete, z.External.Path)

		if z.External.AktePath != "" {
			filesToDelete = append(filesToDelete, z.External.AktePath)
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

package commands

import (
	"flag"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/kilo/store_working_directory"
	"github.com/friedenberg/zit/src/lima/store_with_lock"
	"github.com/friedenberg/zit/src/mike/user_ops"
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

func (c Clean) RunWithLockedStore(
	s store_with_lock.Store,
	args ...string,
) (err error) {
	if len(args) > 0 {
		errors.PrintErrf("args provided will be ignored")
	}

	var possible store_working_directory.CwdFiles

	if possible, err = user_ops.NewGetPossibleZettels(s.Umwelt).Run(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	optionsReadExternal := store_working_directory.OptionsReadExternal{
		Format: zettel.Text{},
	}

	var readResults []zettel_checked_out.Zettel

	readOp := user_ops.ReadCheckedOut{
		Umwelt:              s.Umwelt,
		OptionsReadExternal: optionsReadExternal,
	}

	if readResults, err = readOp.RunMany(s, possible); err != nil {
		err = errors.Wrap(err)
		return
	}

	toDelete := make([]zettel_external.Zettel, 0, len(readResults))
	filesToDelete := make([]string, 0, len(readResults)+len(possible.EmptyDirectories))

	for _, d := range possible.EmptyDirectories {
		filesToDelete = append(filesToDelete, d)
	}

	for _, z := range readResults {
		if z.State != zettel_checked_out.StateExistsAndSame {
			continue
		}

		toDelete = append(toDelete, z.External)

		if z.External.ZettelFD.Path != "" {
			filesToDelete = append(filesToDelete, z.External.ZettelFD.Path)
		}

		if z.External.AkteFD.Path != "" {
			filesToDelete = append(filesToDelete, z.External.AkteFD.Path)
		}
	}

	if s.Umwelt.Konfig.DryRun {
		for _, fOrD := range filesToDelete {
			if pRel, pErr := filepath.Rel(s.Umwelt.Cwd(), fOrD); pErr == nil {
				fOrD = pRel
			}

			stdprinter.Outf("[%s] (would delete)\n", fOrD)
		}

		return
	}

	if err = open_file_guard.DeleteFilesAndDirs(filesToDelete...); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, fOrD := range filesToDelete {
		if pRel, pErr := filepath.Rel(s.Umwelt.Cwd(), fOrD); pErr == nil {
			fOrD = pRel
		}

		stdprinter.Outf("[%s] (deleted)\n", fOrD)
	}

	return
}

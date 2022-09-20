package commands

import (
	"flag"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_external"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/store_objekten"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/kilo/umwelt"
	"github.com/friedenberg/zit/src/lima/user_ops"
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

func (c Clean) Run(
	s *umwelt.Umwelt,
	args ...string,
) (err error) {
	var possible store_working_directory.CwdFiles

	switch {
	case len(args) > 0:
		errors.PrintErrf("Ignoring args")
		fallthrough

	default:
		if possible, err = store_working_directory.MakeCwdFiles(s.Standort().Cwd(), args...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	optionsReadExternal := store_working_directory.OptionsReadExternal{
		Format: zettel.Text{},
	}

	var readResults []zettel_checked_out.Zettel

	readOp := user_ops.ReadCheckedOut{
		Umwelt:              s,
		OptionsReadExternal: optionsReadExternal,
	}

	if readResults, err = readOp.RunMany(possible); err != nil {
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

	for _, ua := range possible.UnsureAkten {
		var szt zettel_transacted.Set

		if szt, err = s.StoreObjekten().ReadAkteSha(ua.Sha); err != nil {
			if errors.Is(err, store_objekten.ErrNotFound{}) {
				continue
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		if szt.Len() > 0 {
			filesToDelete = append(filesToDelete, ua.Path)
		}
	}

	if s.Konfig().DryRun {
		for _, fOrD := range filesToDelete {
			if pRel, pErr := filepath.Rel(s.Standort().Cwd(), fOrD); pErr == nil {
				fOrD = pRel
			}

			errors.PrintOutf("[%s] (would delete)", fOrD)
		}

		return
	}

	if err = files.DeleteFilesAndDirs(filesToDelete...); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, fOrD := range filesToDelete {
		if pRel, pErr := filepath.Rel(s.Standort().Cwd(), fOrD); pErr == nil {
			fOrD = pRel
		}

		errors.PrintOutf("[%s] (deleted)", fOrD)
	}

	return
}

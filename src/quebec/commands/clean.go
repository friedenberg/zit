package commands

import (
	"flag"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/kilo/cwd_files"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/zettel_checked_out"
	"github.com/friedenberg/zit/src/november/store_fs"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/user_ops"
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
	var possible cwd_files.CwdFiles

	switch {
	case len(args) > 0:
		errors.PrintErrf("Ignoring args")
		fallthrough

	default:
		if possible, err = cwd_files.MakeCwdFilesAll(
			s.Konfig(),
			s.Standort().Cwd(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	optionsReadExternal := store_fs.OptionsReadExternal{
		Parser: zettel.MakeTextParser(
			s.StoreObjekten(),
			nil,
		),
	}

	readResults := zettel_checked_out.MakeMutableSetUnique(0)

	readOp := user_ops.ReadCheckedOut{
		Umwelt:              s,
		OptionsReadExternal: optionsReadExternal,
	}

	if err = readOp.RunMany(possible, readResults.Add); err != nil {
		err = errors.Wrap(err)
		return
	}

	toDelete := make([]zettel_external.Zettel, 0, readResults.Len())
	filesToDelete := make([]string, 0, readResults.Len()+len(possible.EmptyDirectories))

	for _, d := range possible.EmptyDirectories {
		filesToDelete = append(filesToDelete, d)
	}

	readResults.Each(
		func(zco *zettel_checked_out.Zettel) (err error) {
			if zco.State != zettel_checked_out.StateExistsAndSame {
				return
			}

			toDelete = append(toDelete, zco.External)

			if zco.External.ZettelFD.Path != "" {
				filesToDelete = append(filesToDelete, zco.External.ZettelFD.Path)
			}

			if zco.External.AkteFD.Path != "" {
				filesToDelete = append(filesToDelete, zco.External.AkteFD.Path)
			}

			return
		},
	)

	//TODO rewrite in verzeichnisseAll
	// for _, ua := range possible.UnsureAkten {
	// 	var szt zettel_transacted.Set

	// 	if szt, err = s.StoreObjekten().ReadAkteSha(ua.Sha); err != nil {
	// 		if errors.Is(err, store_objekten.ErrNotFound{}) {
	// 			continue
	// 		} else {
	// 			err = errors.Wrap(err)
	// 			return
	// 		}
	// 	}

	// 	if szt.Len() > 0 {
	// 		filesToDelete = append(filesToDelete, ua.Path)
	// 	}
	// }

	if s.Konfig().DryRun {
		for _, fOrD := range filesToDelete {
			if pRel, pErr := filepath.Rel(s.Standort().Cwd(), fOrD); pErr == nil {
				fOrD = pRel
			}

			errors.Out().Printf("[%s] (would delete)", fOrD)
		}

		return
	}

	if err = files.DeleteFilesAndDirs(filesToDelete...); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := s.PrinterPathDeleted()

	for _, fOrD := range filesToDelete {
		if pRel, pErr := filepath.Rel(s.Standort().Cwd(), fOrD); pErr == nil {
			fOrD = pRel
		}

		f := &store_fs.Dir{
			Path: fOrD,
		}

		if err = p(f); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

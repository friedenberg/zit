package commands

import (
	"flag"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/juliett/cwd"
	"github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Clean struct{}

func init() {
	registerCommandWithQuery(
		"clean",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Clean{}

			return c
		},
	)
}

func (c Clean) RunWithQuery(
	s *umwelt.Umwelt,
	ms kennung.MetaSet,
) (err error) {
	var possible cwd.CwdFiles

	if possible, err = cwd.MakeCwdFilesMetaSet(
		s.Konfig(),
		s.Standort().Cwd(),
		ms,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	toDelete := make([]objekte.ExternalLike, 0)
	filesToDelete := make([]string, 0)

	for _, d := range possible.EmptyDirectories {
		filesToDelete = append(filesToDelete, d)
	}

	if err = s.StoreWorkingDirectory().ReadFiles(
		possible,
		func(co objekte.CheckedOutLike) (err error) {
			if co.GetState() != objekte.CheckedOutStateExistsAndSame {
				return
			}

			e := co.GetExternal()

			toDelete = append(toDelete, e)

			if ofd := e.GetObjekteFD(); ofd.Path != "" {
				filesToDelete = append(filesToDelete, ofd.Path)
			}

			if afd := e.GetObjekteFD(); afd.Path != "" {
				filesToDelete = append(filesToDelete, afd.Path)
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// optionsReadExternal := store_fs.OptionsReadExternal{}

	// readResults := zettel_checked_out.MakeMutableSetUnique(0)

	// readOp := user_ops.ReadCheckedOut{
	// 	Umwelt:              s,
	// 	OptionsReadExternal: optionsReadExternal,
	// }

	// if err = readOp.RunMany(possible, readResults.Add); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// toDelete := make([]zettel.External, 0, readResults.Len())
	// filesToDelete := make([]string, 0, readResults.Len()+len(possible.EmptyDirectories))

	// readResults.Each(
	// 	func(zco *zettel_checked_out.Zettel) (err error) {
	// 		if zco.State != objekte.CheckedOutStateExistsAndSame {
	// 			return
	// 		}

	// 		toDelete = append(toDelete, zco.External)

	// 		if zco.External.FD.Path != "" {
	// 			filesToDelete = append(filesToDelete, zco.External.FD.Path)
	// 		}

	// 		if zco.External.AkteFD.Path != "" {
	// 			filesToDelete = append(filesToDelete, zco.External.AkteFD.Path)
	// 		}

	// 		return
	// 	},
	// )

	// TODO rewrite in verzeichnisseAll
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

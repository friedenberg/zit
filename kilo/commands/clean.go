package commands

import (
	"flag"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/golf/checkout_store"
	"github.com/friedenberg/zit/india/store_with_lock"
	"github.com/friedenberg/zit/juliett/user_ops"
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
		stdprinter.Errf("args provided will be ignored")
	}

	var possible checkout_store.CwdFiles

	if possible, err = user_ops.NewGetPossibleZettels(s.Umwelt).Run(s); err != nil {
		err = errors.Error(err)
		return
	}

	args = possible.Zettelen

	checkinOptions := checkout_store.CheckinOptions{
		IncludeAkte: true,
		Format:      zettel_formats.Text{},
	}

	var readResults user_ops.ReadCheckedOutResults

	readOp := user_ops.ReadCheckedOut{
		Umwelt:  s.Umwelt,
		Options: checkinOptions,
	}

	if readResults, err = readOp.RunManyStrings(s, args...); err != nil {
		err = errors.Error(err)
		return
	}

	toDelete := make([]stored_zettel.External, 0, len(readResults.Zettelen))
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

	if s.Umwelt.Konfig.DryRun {
		for _, z := range toDelete {
			stdprinter.Outf("[%s] (would delete)\n", z.Hinweis)
		}

		return
	}

	if err = open_file_guard.DeleteFilesAndDirs(filesToDelete...); err != nil {
		err = errors.Error(err)
		return
	}

	for _, z := range toDelete {
		stdprinter.Outf("[%s] (checkout deleted)\n", z.Hinweis)
	}

	return
}

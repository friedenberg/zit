package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/golf/checkout_store"
	"github.com/friedenberg/zit/hotel/zettels"
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

func (c Clean) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) > 0 {
		stdprinter.Errf("args provided will be ignored")
	}

	var possible checkout_store.CwdFiles

	if possible, err = user_ops.NewGetPossibleZettels(u).Run(); err != nil {
		err = errors.Error(err)
		return
	}

	args = possible.Zettelen

	checkinOptions := zettels.CheckinOptions{
		IncludeAkte: true,
		Format:      zettel_formats.Text{},
	}

	var readResults user_ops.ReadCheckedOutResults

	readOp := user_ops.ReadCheckedOut{
		Umwelt:  u,
		Options: checkinOptions,
	}

	if readResults, err = readOp.Run(args...); err != nil {
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

	if u.Konfig.DryRun {
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

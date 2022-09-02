package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/open_file_guard"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/delta/umwelt"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/organize_text"
	"github.com/friedenberg/zit/src/kilo/user_ops"
)

type Add struct {
	Etiketten etikett.Set
	Delete    bool
	OpenAkten bool
	Organize  bool
}

func init() {
	registerCommand(
		"add",
		func(f *flag.FlagSet) Command {
			c := &Add{
				Etiketten: etikett.MakeSet(),
			}

			f.Var(&c.Etiketten, "etiketten", "to add to the created zettels")
			f.BoolVar(&c.Delete, "delete", false, "delete the zettel and akte after successful checkin")
			f.BoolVar(&c.Organize, "organize", false, "")
			f.BoolVar(&c.OpenAkten, "open-akte", false, "also open the Akten")

			return c
		},
	)
}

func (c Add) Run(u *umwelt.Umwelt, args ...string) (err error) {
	zettelsFromAkteOp := user_ops.ZettelFromExternalAkte{
		Umwelt:    u,
		Etiketten: c.Etiketten,
		Delete:    c.Delete,
	}

	var zettelsFromAkteResults zettel_transacted.Set

	if zettelsFromAkteResults, err = zettelsFromAkteOp.Run(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.OpenAkten {
		//TODO
	}

	if !c.Organize {
		return
	}

	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt: u,
		Options: organize_text.Options{
			AssignmentTreeConstructor: organize_text.AssignmentTreeConstructor{
				GroupingEtiketten: etikett.NewSlice(),
				RootEtiketten:     c.Etiketten,
				Transacted:        zettelsFromAkteResults,
			},
		},
	}

	var createOrganizeFileResults user_ops.CreateOrganizeFileResults

	var f *os.File

	if f, err = open_file_guard.TempFileWithPattern("*.md"); err != nil {
		err = errors.Wrap(err)
		return
	}

	if createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(zettelsFromAkteResults, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithFileType("zit-organize").
			Build(),
	}

	if _, err = openVimOp.Run(f.Name()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ot2 organize_text.Text

	readOrganizeTextOp := user_ops.ReadOrganizeFile{}

	if ot2, err = readOrganizeTextOp.RunWithFile(f.Name()); err != nil {
		err = errors.Wrap(err)
		return
	}

	commitOrganizeTextOp := user_ops.CommitOrganizeFile{
		Umwelt: u,
	}

	if _, err = commitOrganizeTextOp.Run(createOrganizeFileResults.Text, ot2); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

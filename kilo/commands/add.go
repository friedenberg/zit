package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/open_file_guard"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/golf/organize_text"
	"github.com/friedenberg/zit/juliett/user_ops"
)

type Add struct {
	Etiketten etikett.Set
	Delete    bool
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

	var zettelsFromAkteResults user_ops.ZettelResults

	if zettelsFromAkteResults, err = zettelsFromAkteOp.Run(args...); err != nil {
		err = errors.Error(err)
		return
	}

	if !c.Organize {
		return
	}

	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt:        u,
		GroupBy:       etikett.NewSlice(),
		RootEtiketten: c.Etiketten,
	}

	var createOrganizeFileResults user_ops.CreateOrganizeFileResults

	var f *os.File

	if f, err = open_file_guard.TempFileWithPattern("*.md"); err != nil {
		err = errors.Error(err)
		return
	}

	if createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(zettelsFromAkteResults, f); err != nil {
		err = errors.Error(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithFileType("zit-organize").
			Build(),
	}

	if _, err = openVimOp.Run(f.Name()); err != nil {
		err = errors.Error(err)
		return
	}

	var ot2 organize_text.Text

	readOrganizeTextOp := user_ops.ReadOrganizeFile{}

	if ot2, err = readOrganizeTextOp.RunWithFile(f.Name()); err != nil {
		err = errors.Error(err)
		return
	}

	commitOrganizeTextOp := user_ops.CommitOrganizeFile{
		Umwelt: u,
	}

	if _, err = commitOrganizeTextOp.Run(createOrganizeFileResults.Text, ot2); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

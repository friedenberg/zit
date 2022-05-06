package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/juliett/user_ops"
)

type Add struct {
	Etiketten _EtikettSet
	Delete    bool
	Organize  bool
}

func init() {
	registerCommand(
		"add",
		func(f *flag.FlagSet) Command {
			c := &Add{
				Etiketten: _EtikettNewSet(),
			}

			f.Var(&c.Etiketten, "etiketten", "to add to the created zettels")
			f.BoolVar(&c.Delete, "delete", false, "delete the zettel and akte after successful checkin")
			f.BoolVar(&c.Organize, "organize", false, "")

			return c
		},
	)
}

func (c Add) Run(u _Umwelt, args ...string) (err error) {
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
		GroupBy:       _EtikettNewSet(),
		RootEtiketten: c.Etiketten,
		GroupByUnique: true,
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
			WithFileType("zit.organize").
			WithSourcedFile("~/.vim/syntax/zit.organize.vim").
			Build(),
	}

	if _, err = openVimOp.Run(f.Name()); err != nil {
		err = errors.Error(err)
		return
	}

	var ot2 _OrganizeText

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

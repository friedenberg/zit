package commands

import (
	"flag"

	"github.com/friedenberg/zit/juliett/user_ops"
)

type Add struct {
	Etiketten _EtikettSet
	Delete    bool
	Organize  bool
	//TODO
	// Edit      bool
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
			// f.BoolVar(&c.Edit, "edit", false, "")

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

	var zettelsFromAkteResults user_ops.ZettelFromExternalAkteResults

	if zettelsFromAkteResults, err = zettelsFromAkteOp.Run(args...); err != nil {
		err = _Error(err)
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

	var createOrganizeFileResults user_ops.CreateOrgaanizeFileResults

	if createOrganizeFileResults, err = createOrganizeFileOp.Run(zettelsFromAkteResults.Zettelen); err != nil {
		err = _Error(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: []string{
			"set ft=zit.organize",
			//TODO find a better solution for this
			"source ~/.vim/syntax/zit.organize.vim",
		},
	}

	if _, err = openVimOp.Run(createOrganizeFileResults.Path); err != nil {
		err = _Error(err)
		return
	}

	var ot2 _OrganizeText

	readOrganizeTextOp := user_ops.ReadOrganizeFile{}

	if ot2, err = readOrganizeTextOp.Run(createOrganizeFileResults.Path); err != nil {
		err = _Error(err)
		return
	}

	commitOrganizeTextOp := user_ops.CommitOrganizeFile{
		Umwelt: u,
	}

	if _, err = commitOrganizeTextOp.Run(createOrganizeFileResults.Text, ot2); err != nil {
		err = _Error(err)
		return
	}

	return
}

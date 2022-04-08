package commands

import (
	"flag"
	"fmt"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/golf/organize_text"
	"github.com/friedenberg/zit/india/store_with_lock"
	"github.com/friedenberg/zit/juliett/user_ops"
)

type Organize struct {
	rootEtiketten _EtikettSet
	GroupBy       _EtikettSet
	GroupByUnique bool
}

func init() {
	registerCommand(
		"organize",
		func(f *flag.FlagSet) Command {
			c := &Organize{
				GroupBy: _EtikettNewSet(),
			}

			f.BoolVar(&c.GroupByUnique, "group-by-unique", false, "group by all unique combinations of etiketten")
			f.Var(&c.GroupBy, "group-by", "etikett prefixes to group zettels")

			return c
		},
	)
}

func (c *Organize) Run(u _Umwelt, args ...string) (err error) {
	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt:        u,
		GroupBy:       c.GroupBy,
		GroupByUnique: c.GroupByUnique,
	}

	if createOrganizeFileOp.RootEtiketten, err = c.getEtikettenFromArgs(args); err != nil {
		err = _Error(err)
		return
	}

	var zettels map[string]_NamedZettel

	if zettels, err = c.getZettels(u, createOrganizeFileOp.RootEtiketten); err != nil {
		err = _Error(err)
		return
	}

	var createOrganizeFileResults user_ops.CreateOrgaanizeFileResults

	if createOrganizeFileResults, err = createOrganizeFileOp.Run(zettels); err != nil {
		err = _Error(err)
		return
	}

OPEN_VIM:
	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithFileType("zit.organize").
			WithSourcedFile("~/.vim/syntax/zit.organize.vim").
			Build(),
	}

	if _, err = openVimOp.Run(createOrganizeFileResults.Path); err != nil {
		err = _Error(err)
		return
	}

	var ot2 _OrganizeText

	readOrganizeTextOp := user_ops.ReadOrganizeFile{}

	if ot2, err = readOrganizeTextOp.Run(createOrganizeFileResults.Path); err != nil {
		if c.handleReadChangesError(err) {
			err = nil
			goto OPEN_VIM
		} else {
			stdprinter.Errf("aborting organize\n")
			return
		}
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

func (c Organize) getEtikettenFromArgs(args []string) (es _EtikettSet, err error) {
	es = _EtikettNewSet()

	for _, s := range args {
		if err = es.AddString(s); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}

func (c Organize) getZettels(u _Umwelt, rootEtiketten _EtikettSet) (zettels map[string]_NamedZettel, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(u); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	if zettels, err = store.Zettels().Query(stored_zettel.FilterEtikettSet(rootEtiketten)); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (c Organize) handleReadChangesError(err error) (tryAgain bool) {
	var errorRead organize_text.ErrorRead

	if err != nil && !errors.As(err, &errorRead) {
		stdprinter.Errf("unrecoverable organize read failure: %s", err)
		tryAgain = false
		return
	}

	stdprinter.Errf("reading changes failed: %q\n", err)
	stdprinter.Errf("would you like to edit and try again? (y/*)\n")

	var answer rune
	var n int

	if n, err = fmt.Scanf("%c", &answer); err != nil {
		tryAgain = false
		stdprinter.Errf("failed to read answer: %s", err)
		return
	}

	if n != 1 {
		tryAgain = false
		stdprinter.Errf("failed to read at exactly 1 answer: %s", err)
		return
	}

	if answer == 'y' || answer == 'Y' {
		tryAgain = true
	}

	return
}

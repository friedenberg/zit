package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/golf/organize_text"
	"github.com/friedenberg/zit/juliett/user_ops"
)

type Organize struct {
	rootEtiketten etikett.Set
	GroupBy       etikett.Set
	GroupByUnique bool
}

func init() {
	registerCommand(
		"organize",
		func(f *flag.FlagSet) Command {
			c := &Organize{
				GroupBy: etikett.NewSet(),
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
		err = errors.Error(err)
		return
	}

	var getResults user_ops.ZettelResults

	getOp := user_ops.GetZettelsFromQuery{Umwelt: u}

	query := stored_zettel.FilterEtikettSet{Set: createOrganizeFileOp.RootEtiketten}

	if getResults, err = getOp.Run(query); err != nil {
		err = errors.Error(err)
		return
	}

	stdoutIsTty := open_file_guard.IsTty(os.Stdout)
	stdinIsTty := open_file_guard.IsTty(os.Stdin)

	if !stdinIsTty && !stdoutIsTty {
		//generate organize, read from stdin,  commit

		createOrganizeFileResults := user_ops.CreateOrganizeFileResults{}

		var f *os.File

		if f, err = open_file_guard.TempFileWithPattern("*.md"); err != nil {
			err = _Error(err)
			return
		}

		if createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(getResults, f); err != nil {
			err = errors.Error(err)
			return
		}

		var ot2 organize_text.Text

		readOrganizeTextOp := user_ops.ReadOrganizeFile{
			Reader: os.Stdin,
		}

		if ot2, err = readOrganizeTextOp.Run(); err != nil {
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
	} else if !stdoutIsTty {
		//generate organize file and write to stdout
		if _, err = createOrganizeFileOp.RunAndWrite(getResults, os.Stdout); err != nil {
			err = errors.Error(err)
			return
		}
	} else {
		//generate temp file, write organize, open vim to edit, commit results
		createOrganizeFileResults := user_ops.CreateOrganizeFileResults{}

		var f *os.File

		if f, err = open_file_guard.TempFileWithPattern("*.md"); err != nil {
			err = _Error(err)
			return
		}

		if createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(getResults, f); err != nil {
			err = errors.Error(err)
			return
		}

		var ot2 organize_text.Text

		if ot2, err = c.readFromVim(f.Name(), createOrganizeFileResults); err != nil {
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
	}

	return
}

func (c Organize) readFromVim(f string, results user_ops.CreateOrganizeFileResults) (ot organize_text.Text, err error) {
	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithFileType("zit.organize").
			WithSourcedFile("~/.vim/syntax/zit.organize.vim").
			Build(),
	}

	if _, err = openVimOp.Run(f); err != nil {
		err = errors.Error(err)
		return
	}

	readOrganizeTextOp := user_ops.ReadOrganizeFile{}

	if ot, err = readOrganizeTextOp.RunWithFile(f); err != nil {
		if c.handleReadChangesError(err) {
			err = nil
			ot, err = c.readFromVim(f, results)
		} else {
			stdprinter.Errf("aborting organize\n")
			return
		}
	}

	return
}

func (c Organize) getEtikettenFromArgs(args []string) (es etikett.Set, err error) {
	es = etikett.NewSet()

	for _, s := range args {
		if err = es.AddString(s); err != nil {
			err = errors.Error(err)
			return
		}
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

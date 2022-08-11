package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/golf/organize_text"
	"github.com/friedenberg/zit/juliett/user_ops"
)

type Organize struct {
	GroupBy        etikett.Slice
	ExtraEtiketten etikett.Set
	Mode           organizeMode
}

type organizeMode int

const (
	organizeModeInteractive = organizeMode(iota)
	organizeModeCommitDirectly
	organizeModeOutputOnly
	organizeModeUnknown = -1
)

func (m *organizeMode) Set(v string) (err error) {
	switch strings.ToLower(v) {
	case "interactive":
		*m = organizeModeInteractive
	case "commit-directly":
		*m = organizeModeCommitDirectly
	case "output-only":
		*m = organizeModeOutputOnly
	default:
		*m = organizeModeUnknown
		err = errors.Errorf("unsupported mode: %s", v)
	}

	return
}

func (m organizeMode) String() string {
	switch m {
	case organizeModeInteractive:
		return "interactive"
	case organizeModeCommitDirectly:
		return "commit-directly"
	case organizeModeOutputOnly:
		return "output-only"
	default:
		return "unknown"
	}
}

func init() {
	registerCommand(
		"organize",
		func(f *flag.FlagSet) Command {
			c := &Organize{
				GroupBy:        etikett.NewSlice(),
				ExtraEtiketten: etikett.MakeSet(),
			}

			f.Var(&c.GroupBy, "group-by", "etikett prefixes to group zettels")
			f.Var(&c.ExtraEtiketten, "extras", "etiketten to always add to the organize text")
			f.Var(&c.Mode, "mode", "mode used for handling stdin and stdout")

			return c
		},
	)
}

func (c *Organize) Run(u *umwelt.Umwelt, args ...string) (err error) {
	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt:         u,
		GroupBy:        c.GroupBy,
		ExtraEtiketten: c.ExtraEtiketten,
		// GroupByUnique: c.GroupByUnique,
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

	switch c.Mode {
	case organizeModeCommitDirectly:
		logz.Print("neither stdin or stdout is a tty")
		logz.Print("generate organize, read from stdin, commit")

		createOrganizeFileResults := user_ops.CreateOrganizeFileResults{}

		var f *os.File

		if f, err = open_file_guard.TempFileWithPattern("*.md"); err != nil {
			err = errors.Error(err)
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

	case organizeModeOutputOnly:
		logz.Print("generate organize file and write to stdout")
		if _, err = createOrganizeFileOp.RunAndWrite(getResults, os.Stdout); err != nil {
			err = errors.Error(err)
			return
		}

	case organizeModeInteractive:
		logz.Print("generate temp file, write organize, open vim to edit, commit results")
		createOrganizeFileResults := user_ops.CreateOrganizeFileResults{}

		var f *os.File

		if f, err = open_file_guard.TempFileWithPattern("*.md"); err != nil {
			err = errors.Error(err)
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

	default:
		err = errors.Errorf("unknown mode")
		return
	}

	return
}

func (c Organize) readFromVim(f string, results user_ops.CreateOrganizeFileResults) (ot organize_text.Text, err error) {
	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithFileType("zit-organize").
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
	es = etikett.MakeSet()

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

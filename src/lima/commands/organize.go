package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/organize_text"
	"github.com/friedenberg/zit/src/juliett/umwelt"
	"github.com/friedenberg/zit/src/kilo/user_ops"
)

type Organize struct {
	organize_text.Options
	Mode organizeMode
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
				Options: organize_text.MakeOptions(),
			}

			f.Var(&c.Mode, "mode", "mode used for handling stdin and stdout")
			c.Options.AddToFlagSet(f)

			return c
		},
	)
}

func (c *Organize) Run(u *umwelt.Umwelt, args ...string) (err error) {
  c.Options.Konfig = u.Konfig()
	c.Options.Abbr = u.StoreObjekten()

	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt:  u,
		Options: c.Options,
	}

	if createOrganizeFileOp.RootEtiketten, err = c.getEtikettenFromArgs(args); err != nil {
		err = errors.Wrap(err)
		return
	}

	var getResults zettel_transacted.Set

	getOp := user_ops.GetZettelsFromQuery{Umwelt: u}

	query := zettel_named.FilterEtikettSet{Set: createOrganizeFileOp.RootEtiketten}

	if getResults, err = getOp.Run(query); err != nil {
		err = errors.Wrap(err)
		return
	}

	createOrganizeFileOp.Transacted = getResults

	switch c.Mode {
	case organizeModeCommitDirectly:
		errors.Print("neither stdin or stdout is a tty")
		errors.Print("generate organize, read from stdin, commit")

		var createOrganizeFileResults *organize_text.Text

		var f *os.File

		if f, err = files.TempFileWithPattern("*.md"); err != nil {
			err = errors.Wrap(err)
			return
		}

		if createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(getResults, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		var ot2 *organize_text.Text

		readOrganizeTextOp := user_ops.ReadOrganizeFile{
			Umwelt: u,
			Reader: os.Stdin,
		}

		if ot2, err = readOrganizeTextOp.Run(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.Lock(); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer u.Unlock()

		commitOrganizeTextOp := user_ops.CommitOrganizeFile{
			Umwelt: u,
		}

		if _, err = commitOrganizeTextOp.Run(createOrganizeFileResults, ot2); err != nil {
			err = errors.Wrap(err)
			return
		}

	case organizeModeOutputOnly:
		errors.Print("generate organize file and write to stdout")
		if _, err = createOrganizeFileOp.RunAndWrite(getResults, os.Stdout); err != nil {
			err = errors.Wrap(err)
			return
		}

	case organizeModeInteractive:
		errors.Print("generate temp file, write organize, open vim to edit, commit results")
		var createOrganizeFileResults *organize_text.Text

		var f *os.File

		if f, err = files.TempFileWithPattern("*.md"); err != nil {
			err = errors.Wrap(err)
			return
		}

		if createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(getResults, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		var ot2 *organize_text.Text

		if ot2, err = c.readFromVim(u, f.Name(), createOrganizeFileResults); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.Lock(); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer u.Unlock()

		commitOrganizeTextOp := user_ops.CommitOrganizeFile{
			Umwelt: u,
		}

		if _, err = commitOrganizeTextOp.Run(createOrganizeFileResults, ot2); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = errors.Errorf("unknown mode")
		return
	}

	return
}

func (c Organize) readFromVim(
	u *umwelt.Umwelt,
	f string,
	results *organize_text.Text,
) (ot *organize_text.Text, err error) {
	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithFileType("zit-organize").
			Build(),
	}

	if _, err = openVimOp.Run(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	readOrganizeTextOp := user_ops.ReadOrganizeFile{
		Umwelt: u,
	}

	if ot, err = readOrganizeTextOp.RunWithFile(f); err != nil {
		if c.handleReadChangesError(err) {
			err = nil
			ot, err = c.readFromVim(u, f, results)
		} else {
			errors.PrintErrf("aborting organize")
			return
		}
	}

	return
}

func (c Organize) getEtikettenFromArgs(args []string) (es etikett.Set, err error) {
	es = etikett.MakeSet()

	for _, s := range args {
		if err = es.AddString(s); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Organize) handleReadChangesError(err error) (tryAgain bool) {
	var errorRead organize_text.ErrorRead

	if err != nil && !errors.As(err, &errorRead) {
		errors.PrintErrf("unrecoverable organize read failure: %s", err)
		tryAgain = false
		return
	}

	errors.PrintErrf("reading changes failed: %q", err)
	errors.PrintErrf("would you like to edit and try again? (y/*)")

	var answer rune
	var n int

	if n, err = fmt.Scanf("%c", &answer); err != nil {
		tryAgain = false
		errors.PrintErrf("failed to read answer: %s", err)
		return
	}

	if n != 1 {
		tryAgain = false
		errors.PrintErrf("failed to read at exactly 1 answer: %s", err)
		return
	}

	if answer == 'y' || answer == 'Y' {
		tryAgain = true
	}

	return
}

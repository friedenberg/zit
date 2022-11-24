package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/id_set"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/organize_text"
	"github.com/friedenberg/zit/src/juliett/zettel_verzeichnisse"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
)

type Organize struct {
	Or bool
	organize_text.Options
	Mode organizeMode

	Filter script_value.ScriptValue
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

			f.BoolVar(&c.Or, "or", false, "allow optional criteria instead of required")
			f.Var(&c.Filter, "filter", "a script to run for each file to transform it the standard zettel format")
			f.Var(&c.Mode, "mode", "mode used for handling stdin and stdout")

			c.Options.AddToFlagSet(f)

			return commandWithIds{
				CommandWithIds: c,
			}
		},
	)
}

func (c Organize) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet(
		id_set.ProtoId{
			MutableId: &kennung.Etikett{},
			Expand: func(v string) (out string, err error) {
				var e kennung.Etikett
				e, err = u.StoreObjekten().ExpandEtikettString(v)
				out = e.String()
				return
			},
		},
		id_set.ProtoId{
			MutableId: &hinweis.Hinweis{},
			Expand: func(v string) (out string, err error) {
				var h hinweis.Hinweis
				h, err = u.StoreObjekten().ExpandHinweisString(v)
				out = h.String()
				return
			},
		},
		id_set.ProtoId{
			MutableId: &kennung.Typ{},
		},
		id_set.ProtoId{
			MutableId: &ts.Time{},
		},
	)

	return
}

func (c *Organize) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	c.Options.Konfig = u.Konfig()
	c.Options.Abbr = u.StoreObjekten()

	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt:  u,
		Options: c.Options,
	}

	createOrganizeFileOp.RootEtiketten = ids.Etiketten()

	typen := ids.Typen()

	switch len(typen) {
	case 0:
		break

	case 1:
		createOrganizeFileOp.Typ = typen[0]

	default:
		err = errors.Errorf("only one typ is supported for organize, but got %q", typen)
		return
	}

	getResults := zettel_transacted.MakeMutableSetUnique(0)

	query := zettel.FilterIdSet{
		Set: ids,
		Or:  c.Or,
	}

	wk := zettel_verzeichnisse.MakeWriterKonfig(u.Konfig())
	w := zettel_verzeichnisse.MakeWriterZettelTransacted(
		collections.MakeChain(
			zettel_transacted.MakeWriterZettelNamed(
				query.WriteZettelNamed,
			),
			getResults.AddAndDoNotRepool,
		),
	)

	if err = u.StoreObjekten().ReadAllSchwanzenVerzeichnisse(wk, w); err != nil {
		err = errors.Wrap(err)
		return
	}

	filterOp := user_ops.FilterZettelsWithScript{
		Set:    getResults,
		Filter: c.Filter,
	}

	if err = filterOp.Run(); err != nil {
		err = errors.Wrap(err)
		return
	}

	createOrganizeFileOp.Transacted = getResults

	switch c.Mode {
	case organizeModeCommitDirectly:
		errors.Log().Print("neither stdin or stdout is a tty")
		errors.Log().Print("generate organize, read from stdin, commit")

		var createOrganizeFileResults *organize_text.Text

		var f *os.File

		if f, err = files.TempFileWithPattern("*." + u.Konfig().DefaultOrganizeExt); err != nil {
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
		errors.Log().Print("generate organize file and write to stdout")
		if _, err = createOrganizeFileOp.RunAndWrite(getResults, os.Stdout); err != nil {
			err = errors.Wrap(err)
			return
		}

	case organizeModeInteractive:
		errors.Log().Print("generate temp file, write organize, open vim to edit, commit results")
		var createOrganizeFileResults *organize_text.Text

		var f *os.File

		if f, err = files.TempFileWithPattern("*." + u.Konfig().DefaultOrganizeExt); err != nil {
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

		if err = u.Reset(); err != nil {
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

	if _, err = openVimOp.Run(u, f); err != nil {
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
			errors.Err().Printf("aborting organize")
			return
		}
	}

	return
}

func (c Organize) handleReadChangesError(err error) (tryAgain bool) {
	var errorRead organize_text.ErrorRead

	if err != nil && !errors.As(err, &errorRead) {
		errors.Err().Printf("unrecoverable organize read failure: %s", err)
		tryAgain = false
		return
	}

	errors.Err().Printf("reading changes failed: %q", err)
	errors.Err().Printf("would you like to edit and try again? (y/*)")

	var answer rune
	var n int

	if n, err = fmt.Scanf("%c", &answer); err != nil {
		tryAgain = false
		errors.Err().Printf("failed to read answer: %s", err)
		return
	}

	if n != 1 {
		tryAgain = false
		errors.Err().Printf("failed to read at exactly 1 answer: %s", err)
		return
	}

	if answer == 'y' || answer == 'Y' {
		tryAgain = true
	}

	return
}

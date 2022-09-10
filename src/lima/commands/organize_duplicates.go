package commands

import (
	"flag"
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/organize_text"
	"github.com/friedenberg/zit/src/juliett/store_with_lock"
	"github.com/friedenberg/zit/src/juliett/umwelt"
	"github.com/friedenberg/zit/src/kilo/user_ops"
)

type OrganizeDuplicates struct {
	organize_text.Options
	Mode organizeMode
}

func init() {
	registerCommand(
		"organize-duplicates",
		func(f *flag.FlagSet) Command {
			c := &OrganizeDuplicates{
				Options: organize_text.Options{
					AssignmentTreeConstructor: organize_text.AssignmentTreeConstructor{
						GroupingEtiketten: etikett.NewSlice(),
						ExtraEtiketten:    etikett.MakeSet(),
					},
				},
			}

			f.Var(&c.GroupingEtiketten, "group-by", "etikett prefixes to group zettels")
			f.Var(&c.ExtraEtiketten, "extras", "etiketten to always add to the organize text")
			f.Var(&c.Mode, "mode", "mode used for handling stdin and stdout")

			return c
		},
	)
}

func (c *OrganizeDuplicates) Run(u *umwelt.Umwelt, args ...string) (err error) {
	// createOrganizeFileOp := user_ops.CreateOrganizeFile{
	// 	Umwelt:  u,
	// 	Options: c.Options,
	// }

	var getResults map[sha.Sha]zettel_transacted.Set

	var store store_with_lock.Store

	if store, err = store_with_lock.New(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	if getResults, err = store.StoreObjekten().ReadAktenDuplicates(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for s, zs := range getResults {
		var zst zettel_transacted.Set

		zst, err = zs.Filter(
			zettel_transacted.MakeSetKeyFuncHinweis(),
			func(zt zettel_transacted.Zettel) (ok bool, err1 error) {
				ok, err1 = store.StoreObjekten().IsSchwanz(zt)
				return
			},
		)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		if zst.Len() <= 1 {
			continue
		}

		errors.PrintOut(s)
		zst.Each(
			func(zt zettel_transacted.Zettel) (err error) {
				errors.PrintOutf("\t%s", zt.Named)
				return
			},
		)
	}

	// createOrganizeFileOp.Transacted = getResults

	// switch c.Mode {
	// case organizeModeCommitDirectly:
	// 	errors.Print("neither stdin or stdout is a tty")
	// 	errors.Print("generate organize, read from stdin, commit")

	// 	createOrganizeFileResults := user_ops.CreateOrganizeFileResults{}

	// 	var f *os.File

	// 	if f, err = files.TempFileWithPattern("*.md"); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	if createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(getResults, f); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	var ot2 organize_text.Text

	// 	readOrganizeTextOp := user_ops.ReadOrganizeFile{
	// 		Reader: os.Stdin,
	// 	}

	// 	if ot2, err = readOrganizeTextOp.Run(); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	commitOrganizeTextOp := user_ops.CommitOrganizeFile{
	// 		Umwelt: u,
	// 	}

	// 	if _, err = commitOrganizeTextOp.Run(createOrganizeFileResults.Text, ot2); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// case organizeModeOutputOnly:
	// 	errors.Print("generate organize file and write to stdout")
	// 	if _, err = createOrganizeFileOp.RunAndWrite(getResults, os.Stdout); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// case organizeModeInteractive:
	// 	errors.Print("generate temp file, write organize, open vim to edit, commit results")
	// 	createOrganizeFileResults := user_ops.CreateOrganizeFileResults{}

	// 	var f *os.File

	// 	if f, err = files.TempFileWithPattern("*.md"); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	if createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(getResults, f); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	var ot2 organize_text.Text

	// 	if ot2, err = c.readFromVim(f.Name(), createOrganizeFileResults); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	commitOrganizeTextOp := user_ops.CommitOrganizeFile{
	// 		Umwelt: u,
	// 	}

	// 	if _, err = commitOrganizeTextOp.Run(createOrganizeFileResults.Text, ot2); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// default:
	// 	err = errors.Errorf("unknown mode")
	// 	return
	// }

	return
}

func (c OrganizeDuplicates) readFromVim(f string, results user_ops.CreateOrganizeFileResults) (ot organize_text.Text, err error) {
	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithFileType("zit-organize").
			Build(),
	}

	if _, err = openVimOp.Run(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	readOrganizeTextOp := user_ops.ReadOrganizeFile{}

	if ot, err = readOrganizeTextOp.RunWithFile(f); err != nil {
		if c.handleReadChangesError(err) {
			err = nil
			ot, err = c.readFromVim(f, results)
		} else {
			errors.PrintErrf("aborting organize")
			return
		}
	}

	return
}

func (c OrganizeDuplicates) getEtikettenFromArgs(args []string) (es etikett.Set, err error) {
	es = etikett.MakeSet()

	for _, s := range args {
		if err = es.AddString(s); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c OrganizeDuplicates) handleReadChangesError(err error) (tryAgain bool) {
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

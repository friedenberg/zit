package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/organize_text"
	"github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
)

type Add struct {
	Dedupe    bool
	Delete    bool
	OpenAkten bool
	Organize  bool
	Filter    script_value.ScriptValue

	zettel.ProtoZettel
}

func init() {
	registerCommandWithQuery(
		"add",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Add{
				ProtoZettel: zettel.MakeEmptyProtoZettel(),
			}

			f.BoolVar(&c.Dedupe, "dedupe", false, "deduplicate added Zettelen based on Akte sha")
			f.BoolVar(&c.Delete, "delete", false, "delete the zettel and akte after successful checkin")
			f.BoolVar(&c.OpenAkten, "open-akten", false, "also open the Akten")
			f.BoolVar(&c.Organize, "organize", false, "")
			c.ProtoZettel.AddToFlagSet(f)

			errors.TodoP2("add support for restricted query to specific gattung")
			return c
		},
	)
}

func (c Add) RunWithQuery(u *umwelt.Umwelt, ms kennung.MetaSet) (err error) {
	zettelsFromAkteOp := user_ops.ZettelFromExternalAkte{
		Umwelt:      u,
		ProtoZettel: c.ProtoZettel,
		Filter:      c.Filter,
		Delete:      c.Delete,
		Dedupe:      c.Dedupe,
	}

	var zettelsFromAkteResults zettel.MutableSet

	if zettelsFromAkteResults, err = zettelsFromAkteOp.Run(ms); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.openAktenIfNecessary(u, zettelsFromAkteResults); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !c.Organize {
		return
	}

	otFlags := organize_text.MakeFlags()
	otFlags.Abbr = u.StoreObjekten().GetAbbrStore().AbbreviateHinweis
	otFlags.RootEtiketten = c.Etiketten
	otFlags.Transacted = zettelsFromAkteResults

	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt:  u,
		Options: otFlags.GetOptions(),
	}

	var createOrganizeFileResults *organize_text.Text

	var f *os.File

	if f, err = files.TempFileWithPattern(
		"*." + u.Konfig().FileExtensions.Organize,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(
		zettelsFromAkteResults,
		f,
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithFileType("zit-organize").
			Build(),
	}

	if _, err = openVimOp.Run(u, f.Name()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ot2 *organize_text.Text

	readOrganizeTextOp := user_ops.ReadOrganizeFile{
		Umwelt: u,
	}

	if ot2, err = readOrganizeTextOp.RunWithFile(f.Name()); err != nil {
		err = errors.Wrap(err)
		return
	}

	commitOrganizeTextOp := user_ops.CommitOrganizeFile{
		Umwelt: u,
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	if _, err = commitOrganizeTextOp.Run(createOrganizeFileResults, ot2); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Add) openAktenIfNecessary(
	u *umwelt.Umwelt,
	zettels zettel.MutableSet,
) (err error) {
	if !c.OpenAkten {
		return
	}

	hs := zettels.ToSliceHinweisen()
	ids := u.MakeIdSet(kennung.MakeMatcherAlways())

	for _, h := range hs {
		ids.Add(h)
	}

	options := store_fs.CheckoutOptions{
		CheckoutMode: objekte.CheckoutModeAkteOnly,
	}

	var checkoutResults zettel.MutableSetCheckedOut

	query := zettel.WriterIds{
		Filter: kennung.Filter{
			Set: ids,
		},
	}

	if checkoutResults, err = u.StoreWorkingDirectory().Checkout(
		options,
		query.WriteZettelTransacted,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	openOp := user_ops.OpenFiles{}

	var filesAkten []string

	if filesAkten, err = zettel.ToSliceFilesAkten(checkoutResults); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = openOp.Run(u, filesAkten...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

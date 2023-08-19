package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/india/objekte_collections"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
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
	registerCommandWithCwdQuery(
		"add",
		func(f *flag.FlagSet) CommandWithCwdQuery {
			c := &Add{
				ProtoZettel: zettel.MakeEmptyProtoZettel(),
			}

			f.BoolVar(
				&c.Dedupe,
				"dedupe",
				false,
				"deduplicate added Zettelen based on Akte sha",
			)
			f.BoolVar(
				&c.Delete,
				"delete",
				false,
				"delete the zettel and akte after successful checkin",
			)
			f.BoolVar(&c.OpenAkten, "open-akten", false, "also open the Akten")
			f.BoolVar(&c.Organize, "organize", false, "")
			c.ProtoZettel.AddToFlagSet(f)

			errors.TodoP2(
				"add support for restricted query to specific gattung",
			)
			return c
		},
	)
}

func (c Add) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
	)
}

func (c Add) RunWithCwdQuery(
	u *umwelt.Umwelt,
	ms kennung.MetaSet,
	pz *cwd.CwdFiles,
) (err error) {
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

	if err = c.openAktenIfNecessary(u, zettelsFromAkteResults, *pz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !c.Organize {
		return
	}

	otFlags := organize_text.MakeFlags()
	u.ApplyToOrganizeOptions(&otFlags.Options)
	// otFlags.Abbr = u.StoreObjekten().GetAbbrStore().AbbreviateHinweis
	otFlags.RootEtiketten = c.Metadatei.Etiketten
	mwk := objekte_collections.MakeMutableSetMetadateiWithKennung()
	zettelsFromAkteResults.Each(
		func(z *sku.TransactedZettel) (err error) {
			return mwk.Add(z.GetSkuLike())
		},
	)
	otFlags.Transacted = mwk

	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt:  u,
		Options: otFlags.GetOptions(u.Konfig().PrintOptions),
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
	cwd cwd.CwdFiles,
) (err error) {
	if !c.OpenAkten {
		return
	}

	hs := collections_value.MakeMutableValueSet[kennung.Hinweis](nil)

	zettels.Each(
		func(z *sku.TransactedZettel) (err error) {
			return hs.Add(z.GetKennung())
		},
	)

	options := store_fs.CheckoutOptions{
		Cwd:          cwd,
		CheckoutMode: checkout_mode.ModeAkteOnly,
	}

	var checkoutResults zettel.MutableSetCheckedOut

	if checkoutResults, err = u.StoreWorkingDirectory().Checkout(
		options,
		func(z *sku.TransactedZettel) (err error) {
			if !hs.Contains(z.GetKennung()) {
				return iter.MakeErrStopIteration()
			}

			return
		},
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

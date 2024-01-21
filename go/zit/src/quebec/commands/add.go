package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/checkout_options"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/india/objekte_collections"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/kilo/organize_text"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/user_ops"
)

type Add struct {
	Dedupe              bool
	Delete              bool
	OpenAkten           bool
	CheckoutAktenAndRun string
	Organize            bool
	Filter              script_value.ScriptValue

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
			f.StringVar(
				&c.CheckoutAktenAndRun,
				"each-akte",
				"",
				"checkout each Akte and run a utility",
			)
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
	ms matcher.Query,
	pz *cwd.CwdFiles,
) (err error) {
	zettelsFromAkteOp := user_ops.ZettelFromExternalAkte{
		Umwelt:      u,
		ProtoZettel: c.ProtoZettel,
		Filter:      c.Filter,
		Delete:      c.Delete,
		Dedupe:      c.Dedupe,
	}

	var zettelsFromAkteResults sku.TransactedMutableSet

	if zettelsFromAkteResults, err = zettelsFromAkteOp.Run(ms); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.openAktenIfNecessary(u, zettelsFromAkteResults, pz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !c.Organize {
		return
	}

	otFlags := organize_text.MakeFlagsWithMetadatei(c.Metadatei)
	u.ApplyToOrganizeOptions(&otFlags.Options)
	// otFlags.Abbr = u.StoreObjekten().GetAbbrStore().AbbreviateHinweis
	mwk := objekte_collections.MakeMutableSetMetadateiWithKennung()
	zettelsFromAkteResults.Each(
		func(z *sku.Transacted) (err error) {
			return mwk.Add(z)
		},
	)
	otFlags.Transacted = mwk

	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt: u,
		Options: otFlags.GetOptions(
			u.Konfig().PrintOptions,
			ms,
			u.SkuFormatOldOrganize(),
			u.SkuFmtNewOrganize(),
		),
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

	if ot2, err = readOrganizeTextOp.RunWithFile(f.Name(), ms); err != nil {
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

	if _, err = commitOrganizeTextOp.Run(
		u,
		createOrganizeFileResults,
		ot2,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Add) openAktenIfNecessary(
	u *umwelt.Umwelt,
	zettels sku.TransactedMutableSet,
	cwd *cwd.CwdFiles,
) (err error) {
	if !c.OpenAkten && c.CheckoutAktenAndRun == "" {
		return
	}

	hs := collections_value.MakeMutableValueSet[values.String](nil)

	zettels.Each(
		func(z *sku.Transacted) (err error) {
			return hs.Add(values.MakeString(z.GetKennung().String()))
		},
	)

	options := checkout_options.Options{
		CheckoutMode: checkout_mode.ModeAkteOnly,
	}

	var checkoutResults sku.CheckedOutMutableSet

	if checkoutResults, err = u.StoreObjekten().Checkout(
		options,
		u.StoreObjekten().MakeReadAllSchwanzen(gattung.Zettel),
		func(z *sku.Transacted) (err error) {
			if !hs.ContainsKey(z.GetKennung().String()) {
				return iter.MakeErrStopIteration()
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var filesAkten []string

	if filesAkten, err = objekte_collections.ToSliceFilesAkten(checkoutResults); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.OpenAkten {
		openOp := user_ops.OpenFiles{}

		if err = openOp.Run(u, filesAkten...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if c.CheckoutAktenAndRun != "" {
		eachAkteOp := user_ops.EachAkte{}

		if err = eachAkteOp.Run(u, c.CheckoutAktenAndRun, filesAkten...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

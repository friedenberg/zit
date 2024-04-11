package commands

import (
	"flag"
	"os"
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/values"
	"code.linenisgreat.com/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/src/charlie/files"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/objekte_collections"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/kilo/cwd"
	"code.linenisgreat.com/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
	"code.linenisgreat.com/zit/src/papa/user_ops"
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
	registerCommandWithQuery(
		"add",
		func(f *flag.FlagSet) CommandWithQuery {
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
			c.AddToFlagSet(f)

			errors.TodoP2(
				"add support for restricted query to specific gattung",
			)
			return c
		},
	)
}

func (c Add) ModifyBuilder(b *query.Builder) {
	b.WithDefaultGattungen(kennung.MakeGattung(gattung.Zettel)).
		WithDoNotMatchEmpty()
}

func (c Add) RunWithQuery(
	u *umwelt.Umwelt,
	qg *query.Group,
) (err error) {
	pz := u.GetStore().GetCwdFiles()

	zettelsFromAkteOp := user_ops.ZettelFromExternalAkte{
		Umwelt:      u,
		ProtoZettel: c.ProtoZettel,
		Filter:      c.Filter,
		Delete:      c.Delete,
		Dedupe:      c.Dedupe,
	}

	var zettelsFromAkteResults sku.TransactedMutableSet

	if zettelsFromAkteResults, err = zettelsFromAkteOp.Run(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.openAktenIfNecessary(
		u,
		qg,
		zettelsFromAkteResults,
		pz,
	); err != nil {
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
			qg,
			u.SkuFormatOldOrganize(),
			u.SkuFmtNewOrganize(),
			u.MakeKennungExpanders(),
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

	if ot2, err = readOrganizeTextOp.RunWithFile(f.Name(), qg); err != nil {
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
		mwk,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Add) openAktenIfNecessary(
	u *umwelt.Umwelt,
	qg *query.Group,
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

	var filesAkten []string
	var l sync.Mutex

	if err = u.GetStore().CheckoutQuery(
		options,
		qg,
		func(z *sku.CheckedOut) (err error) {
			if !hs.ContainsKey(z.Internal.GetKennung().String()) {
				return iter.MakeErrStopIteration()
			}

			e := z.External.GetAkteFD().GetPath()

			if e == "" {
				return iter.MakeErrStopIteration()
			}

			l.Lock()
			filesAkten = append(filesAkten, e)
			l.Unlock()

			return
		},
	); err != nil {
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

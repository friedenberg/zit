package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/charlie/typ"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
	"github.com/friedenberg/zit/src/india/organize_text"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/lima/store_working_directory"
	"github.com/friedenberg/zit/src/mike/umwelt"
	"github.com/friedenberg/zit/src/november/user_ops"
)

type Add struct {
	Dedupe    bool
	Delete    bool
	OpenAkten bool
	Organize  bool
	Filter    script_value.ScriptValue

	//TODO move to protozettel
	Etiketten etikett.Set
	typ.Typ
}

func init() {
	registerCommand(
		"add",
		func(f *flag.FlagSet) Command {
			c := &Add{
				//TODO move to proper place
				Typ: typ.Make("md"),
			}

			f.BoolVar(&c.Dedupe, "dedupe", false, "deduplicate added Zettelen based on Akte sha")
			f.BoolVar(&c.Delete, "delete", false, "delete the zettel and akte after successful checkin")
			f.BoolVar(&c.OpenAkten, "open-akten", false, "also open the Akten")
			f.BoolVar(&c.Organize, "organize", false, "")
			f.Var(&c.Etiketten, "etiketten", "to add to the created zettels")
			f.Var(&c.Filter, "filter", "a script to run for each file to transform it the standard zettel format")
			f.Var(&c.Typ, "typ", "the Typ to use for the newly created Zettelen")

			return c
		},
	)
}

func (c Add) Run(u *umwelt.Umwelt, args ...string) (err error) {
	zettelsFromAkteOp := user_ops.ZettelFromExternalAkte{
		Umwelt: u,
		//TODO add Typ
		ProtoZettel: zettel.ProtoZettel{
			Etiketten: c.Etiketten,
		},
		Filter: c.Filter,
		Delete: c.Delete,
		Dedupe: c.Dedupe,
	}

	var zettelsFromAkteResults zettel_transacted.Set

	if zettelsFromAkteResults, err = zettelsFromAkteOp.Run(args...); err != nil {
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

	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt: u,
		Options: organize_text.Options{
			Abbr:              u.StoreObjekten(),
			GroupingEtiketten: etikett.NewSlice(),
			RootEtiketten:     c.Etiketten,
			Transacted:        zettelsFromAkteResults,
		},
	}

	var createOrganizeFileResults *organize_text.Text

	var f *os.File

	if f, err = files.TempFileWithPattern("*.md"); err != nil {
		err = errors.Wrap(err)
		return
	}

	createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(zettelsFromAkteResults, f)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithFileType("zit-organize").
			Build(),
	}

	if _, err = openVimOp.Run(f.Name()); err != nil {
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

	defer u.Unlock()

	if _, err = commitOrganizeTextOp.Run(createOrganizeFileResults, ot2); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Add) openAktenIfNecessary(
	u *umwelt.Umwelt,
	zettels zettel_transacted.Set,
) (err error) {
	if !c.OpenAkten {
		return
	}

	hs := zettels.ToSliceHinweisen()
	ids := id_set.Make(len(hs))

	for _, h := range hs {
		ids.Add(h)
	}

	options := store_working_directory.CheckoutOptions{
		CheckoutMode: store_working_directory.CheckoutModeAkteOnly,
		Format:       zettel.Text{},
	}

	var checkoutResults zettel_checked_out.Set

	query := zettel_transacted.WriterIds(
		zettel_named.FilterIdSet{
			Set: ids,
		},
	)

	if checkoutResults, err = u.StoreWorkingDirectory().Checkout(options, query); err != nil {
		err = errors.Wrap(err)
		return
	}

	openOp := user_ops.OpenFiles{}

	if err = openOp.Run(checkoutResults.ToSliceFilesAkten()...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

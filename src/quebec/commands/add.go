package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/lima/organize_text"
	"github.com/friedenberg/zit/src/mike/zettel_checked_out"
	"github.com/friedenberg/zit/src/november/store_fs"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/user_ops"
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
	registerCommand(
		"add",
		func(f *flag.FlagSet) Command {
			c := &Add{
				ProtoZettel: zettel.MakeEmptyProtoZettel(),
			}

			f.BoolVar(&c.Dedupe, "dedupe", false, "deduplicate added Zettelen based on Akte sha")
			f.BoolVar(&c.Delete, "delete", false, "delete the zettel and akte after successful checkin")
			f.BoolVar(&c.OpenAkten, "open-akten", false, "also open the Akten")
			f.BoolVar(&c.Organize, "organize", false, "")
			c.ProtoZettel.AddToFlagSet(f)

			return c
		},
	)
}

func (c Add) Run(u *umwelt.Umwelt, args ...string) (err error) {
	zettelsFromAkteOp := user_ops.ZettelFromExternalAkte{
		Umwelt:      u,
		ProtoZettel: c.ProtoZettel,
		Filter:      c.Filter,
		Delete:      c.Delete,
		Dedupe:      c.Dedupe,
	}

	var zettelsFromAkteResults zettel.MutableSet

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

	options := organize_text.MakeOptions()
	options.Abbr = u.StoreObjekten()
	options.RootEtiketten = c.Etiketten
	options.Transacted = zettelsFromAkteResults

	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt:  u,
		Options: options,
	}

	var createOrganizeFileResults *organize_text.Text

	var f *os.File

	if f, err = files.TempFileWithPattern(
		"*." + u.Konfig().FileExtensions.Organize,
	); err != nil {
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
	ids := id_set.Make(len(hs))

	for _, h := range hs {
		ids.Add(h)
	}

	options := store_fs.CheckoutOptions{
		CheckoutMode: store_fs.CheckoutModeAkteOnly,
		Format:       zettel.Text{},
	}

	var checkoutResults zettel_checked_out.MutableSet

	query := zettel.WriterIds{
		Filter: id_set.Filter{
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

	if err = openOp.Run(u, checkoutResults.ToSliceFilesAkten()...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

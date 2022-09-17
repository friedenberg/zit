package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/organize_text"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/kilo/umwelt"
	"github.com/friedenberg/zit/src/lima/user_ops"
)

type Add struct {
	Etiketten etikett.Set
	Delete    bool
	OpenAkten bool
	Organize  bool
}

func init() {
	registerCommand(
		"add",
		func(f *flag.FlagSet) Command {
			c := &Add{
				Etiketten: etikett.MakeSet(),
			}

			f.Var(&c.Etiketten, "etiketten", "to add to the created zettels")
			f.BoolVar(&c.Delete, "delete", false, "delete the zettel and akte after successful checkin")
			f.BoolVar(&c.Organize, "organize", false, "")
			f.BoolVar(&c.OpenAkten, "open-akten", false, "also open the Akten")

			return c
		},
	)
}

func (c Add) Run(u *umwelt.Umwelt, args ...string) (err error) {
	ctx := &errors.Ctx{}

	defer func() {
		err = ctx.Error()
	}()

	zettelsFromAkteOp := user_ops.ZettelFromExternalAkte{
		Umwelt:    u,
		Etiketten: c.Etiketten,
		Delete:    c.Delete,
	}

	var zettelsFromAkteResults zettel_transacted.Set

	if zettelsFromAkteResults = zettelsFromAkteOp.Run(ctx, args...); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	if ctx.Err = c.openAktenIfNecessary(u, zettelsFromAkteResults); !ctx.IsEmpty() {
		ctx.Wrap()
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

	if f, ctx.Err = files.TempFileWithPattern("*.md"); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	createOrganizeFileResults, ctx.Err = createOrganizeFileOp.RunAndWrite(zettelsFromAkteResults, f)

	if !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithFileType("zit-organize").
			Build(),
	}

	if _, ctx.Err = openVimOp.Run(f.Name()); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	if ctx.Err = u.Initialize(); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	var ot2 *organize_text.Text

	readOrganizeTextOp := user_ops.ReadOrganizeFile{
		Umwelt: u,
	}

	if ot2, ctx.Err = readOrganizeTextOp.RunWithFile(f.Name()); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	commitOrganizeTextOp := user_ops.CommitOrganizeFile{
		Umwelt: u,
	}

	if _, ctx.Err = commitOrganizeTextOp.Run(createOrganizeFileResults, ot2); !ctx.IsEmpty() {
		ctx.Wrap()
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

	checkoutOp := user_ops.Checkout{
		Umwelt: u,
		CheckoutOptions: store_working_directory.CheckoutOptions{
			CheckoutMode: store_working_directory.CheckoutModeAkteOnly,
			Format:       zettel.Text{},
		},
	}

	hs := zettels.ToSliceHinweisen()

	var checkoutResults zettel_checked_out.Set

	if checkoutResults, err = checkoutOp.RunManyHinweisen(hs...); err != nil {
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

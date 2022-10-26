package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/bezeichnung"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/charlie/typ"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/kilo/store_working_directory"
	"github.com/friedenberg/zit/src/mike/umwelt"
	"github.com/friedenberg/zit/src/november/user_ops"
)

// TODO move to protozettel
type bez struct {
	bezeichnung.Bezeichnung
	wasSet bool
}

type New struct {
	Edit   bool
	Delete bool
	Count  int
	Filter script_value.ScriptValue

	//TODO move to protozettel
	Bezeichnung bez
	Etiketten   etikett.Set
	typ.Typ
}

func (b *bez) Set(v string) (err error) {
	b.wasSet = true
	return b.Bezeichnung.Set(v)
}

func init() {
	registerCommand(
		"new",
		func(f *flag.FlagSet) Command {
			c := &New{
				//TODO move to proper place
				Typ: typ.Make(""),
			}

			f.BoolVar(&c.Delete, "delete", false, "delete the zettel and akte after successful checkin")
			f.BoolVar(&c.Edit, "edit", true, "create a new empty zettel and open EDITOR or VISUAL for editing and then commit the resulting changes")
			f.IntVar(&c.Count, "count", 1, "when creating new empty zettels, how many to create. otherwise ignored")
			f.Var(&c.Bezeichnung, "bezeichnung", "zettel description (will overwrite existing Bezecihnung")
			f.Var(&c.Etiketten, "etiketten", "comma-separated etiketten (will add to existing Etiketten)")
			f.Var(&c.Filter, "filter", "a script to run for each file to transform it the standard zettel format")
			f.Var(&c.Typ, "typ", "the Typ to use for the newly created Zettelen")

			return c
		},
	)
}

func (c New) ValidateFlagsAndArgs(u *umwelt.Umwelt, args ...string) (err error) {
	if u.Konfig().DryRun && len(args) == 0 {
		err = errors.Errorf("when -dry-run is set, paths to existing zettels must be provided")
		return
	}

	return
}

func (c New) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if err = c.ValidateFlagsAndArgs(u, args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	f := zettel.Text{}
	var zsc zettel_checked_out.Set

	if len(args) == 0 {
		if zsc, err = c.writeNewZettels(u, f); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		var zts zettel_transacted.Set

		if zts, err = c.readExistingFilesAsZettels(u, f, args...); err != nil {
			err = errors.Wrap(err)
			return
		}

		if c.Edit {
			options := store_working_directory.CheckoutOptions{
				CheckoutMode: store_working_directory.CheckoutModeZettelAndAkte,
				Format:       zettel.Text{},
			}

			if zsc, err = u.StoreWorkingDirectory().Checkout(options, zts.WriterFilter()); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if c.Edit {
		if err = c.editZettels(u, zsc); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c New) readExistingFilesAsZettels(
	u *umwelt.Umwelt,
	f zettel.Format,
	args ...string,
) (zts zettel_transacted.Set, err error) {
	opCreateFromPath := user_ops.CreateFromPaths{
		Umwelt: u,
		Format: f,
		Filter: c.Filter,
		Delete: c.Delete,
		ProtoZettel: zettel.ProtoZettel{
			Typ:       c.Typ,
			Etiketten: c.Etiketten,
		},
	}

	if c.Bezeichnung.wasSet {
		opCreateFromPath.ProtoZettel.Bezeichnung = &c.Bezeichnung.Bezeichnung
	}

	if zts, err = opCreateFromPath.Run(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c New) writeNewZettels(
	u *umwelt.Umwelt,
	f zettel.Format,
) (zsc zettel_checked_out.Set, err error) {
	emptyOp := user_ops.WriteNewZettels{
		Umwelt:   u,
		CheckOut: c.Edit,
		CheckoutOptions: store_working_directory.CheckoutOptions{
			Format: f,
		},
	}

	var defaultEtiketten etikett.Set

	if defaultEtiketten, err = u.DefaultEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	mes := c.Etiketten.MutableCopy()
	mes.Merge(defaultEtiketten)
	c.Etiketten = mes.Copy()

	z := zettel.Zettel{
		Bezeichnung: c.Bezeichnung.Bezeichnung,
		Etiketten:   c.Etiketten,
		Typ:         c.Typ,
	}

	if zsc, err = emptyOp.RunMany(z, c.Count); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c New) editZettels(
	u *umwelt.Umwelt,
	zsc zettel_checked_out.Set,
) (err error) {
	if !c.Edit {
		errors.Print("edit set to false, not editing")
		return
	}

	fs := zsc.ToSliceFilesZettelen()

	var cwdFiles store_working_directory.CwdFiles

	if cwdFiles, err = store_working_directory.MakeCwdFilesExactly(u.Konfig().Compiled, u.Standort().Cwd(), fs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithCursorLocation(2, 3).
			WithFileType("zit-zettel").
			WithInsertMode().
			Build(),
	}

	if _, err = openVimOp.Run(cwdFiles.ZettelFiles()...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	readOp := user_ops.ReadCheckedOut{
		Umwelt: u,
		OptionsReadExternal: store_working_directory.OptionsReadExternal{
			Format: zettel.Text{},
		},
	}

	var zslc zettel_checked_out.Set

	if zslc, err = readOp.RunMany(cwdFiles); err != nil {
		err = errors.Wrap(err)
		return
	}

	checkinOp := user_ops.Checkin{
		Umwelt:              u,
		OptionsReadExternal: readOp.OptionsReadExternal,
	}

	zsle := zslc.ToSliceZettelsExternal()

	if _, err = checkinOp.Run(zsle...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

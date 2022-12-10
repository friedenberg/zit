package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/cwd_files"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
)

type New struct {
	Edit   bool
	Delete bool
	Dedupe bool
	Count  int
	Filter script_value.ScriptValue

	zettel.ProtoZettel
}

func init() {
	registerCommand(
		"new",
		func(f *flag.FlagSet) Command {
			c := &New{
				ProtoZettel: zettel.MakeProtoZettel(),
			}

			f.BoolVar(&c.Delete, "delete", false, "delete the zettel and akte after successful checkin")
			f.BoolVar(&c.Dedupe, "dedupe", false, "deduplicate added Zettelen based on Akte sha")
			f.BoolVar(&c.Edit, "edit", true, "create a new empty zettel and open EDITOR or VISUAL for editing and then commit the resulting changes")
			f.IntVar(&c.Count, "count", 1, "when creating new empty zettels, how many to create. otherwise ignored")

			f.Var(&c.Filter, "filter", "a script to run for each file to transform it the standard zettel format")
			c.ProtoZettel.AddToFlagSet(f)

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
	var zsc zettel_checked_out.MutableSet

	if len(args) == 0 {
		if zsc, err = c.writeNewZettels(u, f); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		var zts zettel.MutableSet

		if zts, err = c.readExistingFilesAsZettels(u, f, args...); err != nil {
			err = errors.Wrap(err)
			return
		}

		if c.Edit {
			options := store_fs.CheckoutOptions{
				CheckoutMode: store_fs.CheckoutModeZettelAndAkte,
				Format:       zettel.Text{},
			}

			if zsc, err = u.StoreWorkingDirectory().Checkout(
				options,
				zts.WriterContainer(),
			); err != nil {
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
) (zts zettel.MutableSet, err error) {
	opCreateFromPath := user_ops.CreateFromPaths{
		Umwelt:      u,
		Format:      f,
		Filter:      c.Filter,
		Delete:      c.Delete,
		Dedupe:      c.Dedupe,
		ProtoZettel: c.ProtoZettel,
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
) (zsc zettel_checked_out.MutableSet, err error) {
	emptyOp := user_ops.WriteNewZettels{
		Umwelt:   u,
		CheckOut: c.Edit,
		CheckoutOptions: store_fs.CheckoutOptions{
			Format: f,
		},
	}

	var defaultEtiketten kennung.EtikettSet

	if defaultEtiketten, err = u.DefaultEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	mes := c.Etiketten.MutableCopy()
	defaultEtiketten.Each(mes.Add)
	c.Etiketten = mes.Copy()

	if zsc, err = emptyOp.RunMany(c.ProtoZettel, c.Count); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c New) editZettels(
	u *umwelt.Umwelt,
	zsc zettel_checked_out.MutableSet,
) (err error) {
	if !c.Edit {
		errors.Log().Print("edit set to false, not editing")
		return
	}

	fs := zsc.ToSliceFilesZettelen()

	var cwdFiles cwd_files.CwdFiles

	if cwdFiles, err = cwd_files.MakeCwdFilesExactly(u.Konfig(), u.Standort().Cwd(), fs...); err != nil {
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

	if _, err = openVimOp.Run(u, cwdFiles.ZettelFiles()...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	readOp := user_ops.ReadCheckedOut{
		Umwelt: u,
		OptionsReadExternal: store_fs.OptionsReadExternal{
			Format: zettel.Text{},
		},
	}

	zslc := zettel_checked_out.MakeMutableSetUnique(0)

	if err = readOp.RunMany(cwdFiles, zslc.Add); err != nil {
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

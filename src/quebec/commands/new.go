package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/juliett/cwd"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/lima/zettel_checked_out"
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
				ProtoZettel: zettel.MakeEmptyProtoZettel(),
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

	f := zettel.MakeObjekteTextFormat(
		u.StoreObjekten(),
		nil,
	)

	var zsc zettel_checked_out.MutableSet

	if len(args) == 0 {
		if zsc, err = c.writeNewZettels(u, f); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		var zts schnittstellen.MutableSet[*zettel.Transacted]

		if zts, err = c.readExistingFilesAsZettels(u, f, args...); err != nil {
			err = errors.Wrap(err)
			return
		}

		if c.Edit {
			options := store_fs.CheckoutOptions{
				CheckoutMode: objekte.CheckoutModeObjekteAndAkte,
			}

			if zsc, err = u.StoreWorkingDirectory().Checkout(
				options,
				collections.WriterContainer[*zettel.Transacted](zts, collections.MakeErrStopIteration()),
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
	f zettel.ObjekteParser,
	args ...string,
) (zts schnittstellen.MutableSet[*zettel.Transacted], err error) {
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
	f zettel.ObjekteFormatter,
) (zsc zettel_checked_out.MutableSet, err error) {
	emptyOp := user_ops.WriteNewZettels{
		Umwelt:          u,
		CheckOut:        c.Edit,
		CheckoutOptions: store_fs.CheckoutOptions{},
	}

	var defaultEtiketten kennung.EtikettSet

	if defaultEtiketten, err = u.DefaultEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	mes := c.Etiketten.MutableClone()
	defaultEtiketten.Each(mes.Add)
	c.Etiketten = mes.ImmutableClone()

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

	fs := zettel_checked_out.ToSliceFilesZettelen(zsc)

	var cwdFiles cwd.CwdFiles

	if cwdFiles, err = cwd.MakeCwdFilesExactly(u.Konfig(), u.Standort().Cwd(), fs...); err != nil {
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

	checkinOp := user_ops.Checkin{}

	if err = checkinOp.Run(u, cwdFiles); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

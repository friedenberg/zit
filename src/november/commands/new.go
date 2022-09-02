package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/bezeichnung"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/typ"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/delta/umwelt"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/hotel/zettel_external"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/kilo/store_working_directory"
	"github.com/friedenberg/zit/src/lima/store_with_lock"
	"github.com/friedenberg/zit/src/mike/user_ops"
)

type bez struct {
	bezeichnung.Bezeichnung
	wasSet bool
}

type New struct {
	Bezeichnung bez
	Edit        bool
	Delete      bool
	Count       int
	Etiketten   etikett.Set
	Filter      script_value.ScriptValue
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
				Etiketten: etikett.MakeSet(),
			}

			f.Var(&c.Filter, "filter", "a script to run for each file to transform it the standard zettel format")
			f.IntVar(&c.Count, "count", 1, "when creating new empty zettels, how many to create. otherwise ignored")
			f.BoolVar(&c.Edit, "edit", true, "create a new empty zettel and open EDITOR or VISUAL for editing and then commit the resulting changes")
			f.BoolVar(&c.Delete, "delete", false, "delete the zettel and akte after successful checkin")
			f.Var(&c.Bezeichnung, "bezeichnung", "zettel description (will overwrite existing Bezecihnung")
			f.Var(&c.Etiketten, "etiketten", "comma-separated etiketten (will add to existing Etiketten)")

			return c
		},
	)
}

func (c New) ValidateFlagsAndArgs(u *umwelt.Umwelt, args ...string) (err error) {
	if u.Konfig.DryRun && len(args) == 0 {
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
		if zsc, err = c.readExistingFilesAsZettels(u, f, args...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = c.editZettelsIfRequested(u, zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c New) readExistingFilesAsZettels(
	u *umwelt.Umwelt,
	f zettel.Format,
	args ...string,
) (zsc zettel_checked_out.Set, err error) {
	opCreateFromPath := user_ops.CreateFromPaths{
		Umwelt: u,
		Format: f,
		Filter: c.Filter,
		Delete: c.Delete,
	}

	//TODO add bezeichnung and etiketten
	if _, err = opCreateFromPath.Run(args...); err != nil {
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
		Umwelt: u,
		CheckoutOptions: store_working_directory.CheckoutOptions{
			Format: f,
		},
	}

	var defaultEtiketten etikett.Set

	if defaultEtiketten, err = u.DefaultEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.Etiketten.Merge(defaultEtiketten)

	z := zettel.Zettel{
		Bezeichnung: c.Bezeichnung.Bezeichnung,
		Etiketten:   c.Etiketten,
		Typ:         typ.Typ{Value: "md"},
	}

	var s store_with_lock.Store

	if s, err = store_with_lock.New(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.PanicIfError(s.Flush)

	if zsc, err = emptyOp.RunMany(s, z, c.Count); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c New) editZettelsIfRequested(
	u *umwelt.Umwelt,
	zsc zettel_checked_out.Set,
) (err error) {
	if !c.Edit {
		errors.Print("edit set to false, not editing")
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithCursorLocation(2, 3).
			WithInsertMode().
			Build(),
	}

	cwdFiles := store_working_directory.CwdFiles{
		Zettelen: make([]string, 0, zsc.Len()),
	}

	zsc.Each(
		func(zc zettel_checked_out.Zettel) (err error) {
			cwdFiles.Zettelen = append(cwdFiles.Zettelen, zc.External.ZettelFD.Path)
			return nil
		},
	)

	if _, err = openVimOp.Run(cwdFiles.Zettelen...); err != nil {
		err = errors.Wrap(err)
		return
	}

	readOp := user_ops.ReadCheckedOut{
		Umwelt: u,
		OptionsReadExternal: store_working_directory.OptionsReadExternal{
			Format: zettel.Text{},
		},
	}

	var s store_with_lock.Store

	if s, err = store_with_lock.New(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.PanicIfError(s.Flush)

	var zslc []zettel_checked_out.Zettel

	if zslc, err = readOp.RunMany(s, cwdFiles); err != nil {
		err = errors.Wrap(err)
		return
	}

	checkinOp := user_ops.Checkin{
		Umwelt:              s.Umwelt,
		OptionsReadExternal: readOp.OptionsReadExternal,
	}

	zsle := make([]zettel_external.Zettel, len(zslc))

	for i, zc := range zslc {
		zsle[i] = zc.External
	}

	if _, err = checkinOp.Run(s, zsle...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

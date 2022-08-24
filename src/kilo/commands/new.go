package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/bezeichnung"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/typ"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/script_value"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	store_checkout "github.com/friedenberg/zit/src/hotel/store_checkout"
	"github.com/friedenberg/zit/src/india/store_with_lock"
	"github.com/friedenberg/zit/src/juliett/user_ops"
)

type bez struct {
	bezeichnung.Bezeichnung
	wasSet bool
}

type New struct {
	Bezeichnung bez
	Edit        bool
	Delete      bool
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
			f.BoolVar(&c.Edit, "edit", false, "create a new empty zettel and open EDITOR or VISUAL for editing and then commit the resulting changes")
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

	if len(args) > 0 && c.Edit {
		err = errors.Errorf("editing not supported when importing existing zettels")
		return
	}

	return
}

func (c New) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if err = c.ValidateFlagsAndArgs(u, args...); err != nil {
		err = errors.Error(err)
		return
	}

	f := zettel_formats.Text{}

	if len(args) == 0 {
		var cz stored_zettel.CheckedOut

		if cz, err = c.writeNewZettel(u, f); err != nil {
			err = errors.Error(err)
			return
		}

		if err = c.editZettelIfRequested(u, cz); err != nil {
			err = errors.Error(err)
			return
		}
	} else {
		if err = c.readExistingFilesAsZettels(u, f, args...); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}

func (c New) readExistingFilesAsZettels(u *umwelt.Umwelt, f zettel.Format, args ...string) (err error) {
	opCreateFromPath := user_ops.CreateFromPaths{
		Umwelt: u,
		Format: f,
		Filter: c.Filter,
		Delete: c.Delete,
	}

	//TODO add bezeichnung and etiketten
	if _, err = opCreateFromPath.Run(args...); err != nil {
		err = errors.Error(err)
		return
	}

	//TODO if edit, checkout zettels and made editable

	return
}

func (c New) writeNewZettel(
	u *umwelt.Umwelt,
	f zettel.Format,
) (cz stored_zettel.CheckedOut, err error) {
	emptyOp := user_ops.WriteNewZettels{
		Umwelt: u,
		CheckoutOptions: store_checkout.CheckoutOptions{
			Format: f,
		},
	}

	var defaultEtiketten etikett.Set

	if defaultEtiketten, err = u.DefaultEtiketten(); err != nil {
		err = errors.Error(err)
		return
	}

	c.Etiketten.Merge(defaultEtiketten)

	z := zettel.Zettel{
		Bezeichnung: c.Bezeichnung.Bezeichnung,
		Etiketten:   c.Etiketten,
		AkteExt:     akte_ext.AkteExt{Value: "md"},
	}

	var s store_with_lock.Store

	if s, err = store_with_lock.New(u); err != nil {
		err = errors.Error(err)
		return
	}

	defer s.Flush()

	if cz, err = emptyOp.RunOne(s, z); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (c New) editZettelIfRequested(
	u *umwelt.Umwelt,
	cz stored_zettel.CheckedOut,
) (err error) {
	if !c.Edit {
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithCursorLocation(2, 3).
			WithInsertMode().
			Build(),
	}

	if _, err = openVimOp.Run(cz.External.Path); err != nil {
		err = errors.Error(err)
		return
	}

	readOp := user_ops.ReadCheckedOut{
		Umwelt: u,
		OptionsReadExternal: store_checkout.OptionsReadExternal{
			Format: zettel_formats.Text{},
		},
	}

	var s store_with_lock.Store

	if s, err = store_with_lock.New(u); err != nil {
		err = errors.Error(err)
		return
	}

	defer s.Flush()

	if cz, err = readOp.RunOneString(s, cz.External.Path); err != nil {
		err = errors.Error(err)
		return
	}

	checkinOp := user_ops.Checkin{
		Umwelt:              s.Umwelt,
		OptionsReadExternal: readOp.OptionsReadExternal,
	}

	if _, err = checkinOp.Run(s, cz.External); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

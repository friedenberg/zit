package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/bezeichnung"
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/bravo/akte_ext"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/juliett/user_ops"
)

type bez struct {
	bezeichnung.Bezeichnung
	wasSet bool
}

type New struct {
	Bezeichnung bez
	Edit        bool
	Etiketten   etikett.Set
	Filter      _ScriptValue
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
				Etiketten: etikett.NewSet(),
			}

			f.Var(&c.Filter, "filter", "a script to run for each file to transform it the standard zettel format")
			f.BoolVar(&c.Edit, "edit", false, "create a new empty zettel and open EDITOR or VISUAL for editing and then commit the resulting changes")
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

func (c New) Run(u _Umwelt, args ...string) (err error) {
	if err = c.ValidateFlagsAndArgs(u, args...); err != nil {
		err = errors.Error(err)
		return
	}

	f := _ZettelFormatsText{}

	if len(args) == 0 {
		if err = c.writeNewZettel(u, f); err != nil {
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
	}

	//TODO add bezeichnung and etiketten
	if _, err = opCreateFromPath.Run(args...); err != nil {
		err = errors.Error(err)
		return
	}

	//TODO if edit, checkout zettels and made editable

	return
}

func (c New) writeNewZettel(u *umwelt.Umwelt, f zettel.Format) (err error) {
	opCreateFromPath := user_ops.CreateFromPaths{
		Umwelt: u,
		Format: f,
		Filter: c.Filter,
	}

	emptyOp := user_ops.WriteNewZettels{
		Umwelt: u,
		Format: f,
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

	var results stored_zettel.SetExternal

	if results, err = emptyOp.Run(z); err != nil {
		err = errors.Error(err)
		return
	}

	opCreateFromPath.ReadHinweisFromPath = true

	if c.Edit {
		openVimOp := user_ops.OpenVim{
			Options: vim_cli_options_builder.New().
				WithCursorLocation(2, 3).
				WithInsertMode().
				WithFileType("zit.zettel").
				WithSourcedFile("~/.vim/syntax/zit.zettel.vim").
				Build(),
		}

		if _, err = openVimOp.Run(results.Paths()...); err != nil {
			err = errors.Error(err)
			return
		}
	}

	if _, err = opCreateFromPath.Run(results.Paths()...); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

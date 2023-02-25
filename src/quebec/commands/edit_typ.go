package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
)

type EditTyp struct {
	All bool
}

func init() {
	registerCommand(
		"edit-typ",
		func(f *flag.FlagSet) Command {
			c := &EditTyp{}

			f.BoolVar(&c.All, "all", false, "edit all Typen")

			return commandWithIds{CommandWithIds: c}
		},
	)
}

func (c EditTyp) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Typ,
	)
}

func (c EditTyp) RunWithIds(u *umwelt.Umwelt, ids kennung.Set) (err error) {
	tks := ids.Typen.ImmutableClone()

	switch {
	case tks.Len() == 0 && !c.All:
		err = errors.Normalf("No Typen specified for editing. To edit all, use -all.")
		return

	case c.All && tks.Len() > 0:
		errors.Err().Print("Ignoring arguments because -all is set.")

		fallthrough

	case c.All:
		mtks := collections.MakeMutableSetStringer[kennung.Typ]()

		u.Konfig().Typen.Each(
			func(tt *typ.Transacted) (err error) {
				return mtks.Add(tt.Sku.Kennung)
			},
		)

		tks = mtks.ImmutableClone()
	}

	var ps []string

	if ps, err = u.StoreWorkingDirectory().MakeTempTypFiles(tks); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithCursorLocation(2, 3).
			WithFileType("zit-typ").
			WithInsertMode().
			Build(),
	}

	if _, err = openVimOp.Run(u, ps...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var tes schnittstellen.Set[*typ.External]

	if tes, err = c.readTempTypFiles(u, ps); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	if err = tes.Each(
		func(e *typ.External) (err error) {
			var tt *typ.Transacted

			if tt, err = u.StoreObjekten().Typ().CreateOrUpdate(
				&e.Objekte,
				&e.Sku.Kennung,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			u.KonfigPtr().AddTyp(tt)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c EditTyp) readTempTypFiles(
	u *umwelt.Umwelt,
	ps []string,
) (out schnittstellen.Set[*typ.External], err error) {
	ts := collections.MakeMutableSet[*typ.External](
		typ.ExternalKeyer{}.Key,
	)

	formatText := typ.MakeFormatText(u.StoreObjekten())

	iter := func(p string) (err error) {
		var f *os.File

		if f, err = files.Open(p); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, f.Close)

		var fdee kennung.FD

		if fdee, err = kennung.File(f); err != nil {
			err = errors.Wrap(err)
			return
		}

		te := &typ.External{
			Sku: sku.External[kennung.Typ, *kennung.Typ]{
				Kennung: kennung.MustTyp(fdee.FileNameSansExt()),
				FDs: sku.ExternalFDs{
					Objekte: fdee,
				},
			},
		}

		// TODO offer option to edit again
		if _, err = formatText.ReadFormat(f, &te.Objekte); err != nil {
			err = errors.Wrap(err)
			return
		}

		ts.Add(te)

		return
	}

	for _, p := range ps {
		if err = iter(p); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	out = ts.ImmutableClone()

	return
}

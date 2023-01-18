package commands

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/fd"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/india/typ"
	"github.com/friedenberg/zit/src/mike/store_objekten"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/user_ops"
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

func (c EditTyp) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet(
		id_set.ProtoId{
			Setter: &kennung.Typ{},
		},
	)

	return
}

func (c EditTyp) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	tks := ids.Typen.Copy()

	switch {
	case tks.Len() == 0 && !c.All:
		err = errors.Normalf("No Typen specified for editing. To edit all, use -all.")
		return

	case c.All && tks.Len() > 0:
		errors.Err().Print("Ignoring arguments because -all is set.")

		fallthrough

	case c.All:
		mtks := collections.MakeMutableValueSet[kennung.Typ, *kennung.Typ]()

		u.Konfig().Typen.Each(
			func(tt *typ.Transacted) (err error) {
				return mtks.Add(tt.Sku.Kennung)
			},
		)

		tks = mtks.Copy()
	}

	var ps []string

	if ps, err = c.makeTempTypFiles(u, tks); err != nil {
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

	var tes collections.Set[*typ.External]

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

func (c EditTyp) makeTempTypFiles(
	u *umwelt.Umwelt,
	tks collections.ValueSet[kennung.Typ, *kennung.Typ],
) (ps []string, err error) {
	ps = make([]string, 0, tks.Len())

	var tempDir string

	if tempDir, err = files.TempDir(); err != nil {
		err = errors.Wrap(err)
		return
	}

	format := typ.MakeFormatText(u.StoreObjekten())

	if err = tks.Each(
		func(tk kennung.Typ) (err error) {
			var tt *typ.Transacted

			if tt, err = u.StoreObjekten().Typ().ReadOne(&tk); err != nil {
				if errors.Is(err, store_objekten.ErrNotFound{}) {
					err = nil
					tt = &typ.Transacted{
						Sku: sku.Transacted[kennung.Typ, *kennung.Typ]{
							Kennung: tk,
						},
					}
				} else {
					err = errors.Wrap(err)
					return
				}
			}

			var f *os.File

			if f, err = files.CreateExclusiveWriteOnly(
				path.Join(tempDir, fmt.Sprintf("%s.%s", tk.String(), u.Konfig().FileExtensions.Typ)),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.Deferred(&err, f.Close)

			ps = append(ps, f.Name())

			if _, err = format.WriteFormat(f, &tt.Objekte); err != nil {
				err = errors.Wrap(err)
				return
			}

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
) (out collections.Set[*typ.External], err error) {
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

		var fdee fd.FD

		if fdee, err = fd.File(f); err != nil {
			err = errors.Wrap(err)
			return
		}

		te := &typ.External{
			Sku: sku.External[kennung.Typ, *kennung.Typ]{
				Kennung: kennung.MustTyp(fdee.FileNameSansExt()),
			},
			FD: fdee,
		}

		//TODO offer option to edit again
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

	out = ts.Copy()

	return
}

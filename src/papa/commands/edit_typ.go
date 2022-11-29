package commands

import (
	"flag"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/typ_toml"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/echo/id_set"
	"github.com/friedenberg/zit/src/golf/typ"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
)

type EditTyp struct {
}

func init() {
	registerCommand(
		"edit-typ",
		func(f *flag.FlagSet) Command {
			c := &EditTyp{}

			return commandWithIds{CommandWithIds: c}
		},
	)
}

func (c EditTyp) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet(
		id_set.ProtoId{
			MutableId: &kennung.Typ{},
		},
	)

	return
}

func (c EditTyp) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	var ps []string

	if ps, err = c.makeTempTypFiles(u, ids.Typen()); err != nil {
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
			if _, err = u.StoreObjekten().Typ().Update(
				&e.Objekte,
				&e.Kennung,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			u.KonfigPtr().Transacted.Objekte.AddTyp(
				&e.Objekte.Akte,
				&e.Kennung,
			)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	k := u.KonfigPtr()

	if err = k.Transacted.Objekte.Recompile(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = u.StoreObjekten().Konfig().Update(
		&k.Transacted.Objekte,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c EditTyp) makeTempTypFiles(
	u *umwelt.Umwelt,
	tks []kennung.Typ,
) (ps []string, err error) {
	ps = make([]string, 0, len(tks))

	var tempDir string

	if tempDir, err = files.TempDir(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, tk := range tks {
		var tt *typ.Transacted

		if tt, err = u.StoreObjekten().Typ().ReadOne(&tk); err != nil {
			err = errors.Wrap(err)
			return
		}

		format := typ_toml.MakeFormatText(u.StoreObjekten())

		func() {
			var f *os.File

			if f, err = files.CreateExclusiveWriteOnly(
				path.Join(tempDir, tk.String()),
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
		}()
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

	formatText := typ_toml.MakeFormatText(u.StoreObjekten())

	for _, p := range ps {
		func() {
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
				Kennung: kennung.MustTyp(fdee.FileNameSansExt()),
				FD:      fdee,
			}

			//TODO offer option to edit again
			if _, err = formatText.ReadFormat(f, &te.Objekte); err != nil {
				err = errors.Wrap(err)
				return
			}

			ts.Add(te)
		}()
	}

	out = ts.Copy()

	return
}

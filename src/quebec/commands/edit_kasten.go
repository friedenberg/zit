package commands

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/fd"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/kasten"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
)

type EditKasten struct {
}

func init() {
	registerCommand(
		"edit-kasten",
		func(f *flag.FlagSet) Command {
			c := &EditKasten{}

			return commandWithIds{CommandWithIds: c}
		},
	)
}

func (c EditKasten) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Kasten,
	)
}

func (c EditKasten) RunWithIds(u *umwelt.Umwelt, ids kennung.Set) (err error) {
	tks := ids.Kisten.Copy()

	switch {
	case tks.Len() == 0 && !ids.Sigil.IncludesAll():
		err = errors.Normalf("No Kisten specified for editing. To edit all, use -all.")
		return

	case ids.Sigil.IncludesAll() && tks.Len() > 0:
		errors.Err().Print("Ignoring arguments because -all is set.")

		fallthrough

	case ids.Sigil.IncludesAll():
		mtks := collections.MakeMutableValueSet[kennung.Kasten, *kennung.Kasten]()

		u.Konfig().Kisten.Each(
			func(tt *kasten.Transacted) (err error) {
				return mtks.Add(tt.Sku.Kennung)
			},
		)

		tks = mtks.Copy()
	}

	var ps []string

	if ps, err = c.makeTempKastenFiles(u, tks); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithCursorLocation(2, 3).
			WithFileType("zit-kasten").
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

	var tes collections.Set[*kasten.External]

	if tes, err = c.readTempKastenFiles(u, ps); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	if err = tes.Each(
		func(e *kasten.External) (err error) {
			if _, err = u.StoreObjekten().Kasten().CreateOrUpdate(
				&e.Objekte,
				&e.Sku.Kennung,
			); err != nil {
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

func (c EditKasten) makeTempKastenFiles(
	u *umwelt.Umwelt,
	tks collections.ValueSet[kennung.Kasten, *kennung.Kasten],
) (ps []string, err error) {
	ps = make([]string, 0, tks.Len())

	var tempDir string

	if tempDir, err = u.Standort().DirTempOS(); err != nil {
		err = errors.Wrap(err)
		return
	}

	format := kasten.MakeFormatText(u.StoreObjekten())

	if err = tks.Each(
		func(tk kennung.Kasten) (err error) {
			var tt *kasten.Transacted

			if tt, err = u.StoreObjekten().Kasten().ReadOne(&tk); err != nil {
				if errors.Is(err, objekte_store.ErrNotFound{}) {
					err = nil
					tt = &kasten.Transacted{
						Sku: sku.Transacted[kennung.Kasten, *kennung.Kasten]{
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
				path.Join(
					tempDir,
					fmt.Sprintf("%s.%s", tk.String(), u.Konfig().FileExtensions.Kasten),
				),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.Deferred(&err, f.Close)

			ps = append(ps, f.Name())

			if _, err = format.Format(f, &tt.Objekte); err != nil {
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

func (c EditKasten) readTempKastenFiles(
	u *umwelt.Umwelt,
	ps []string,
) (out collections.Set[*kasten.External], err error) {
	ts := collections.MakeMutableSet[*kasten.External](
		kasten.ExternalKeyer{}.Key,
	)

	formatText := kasten.MakeFormatText(u.StoreObjekten())

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

			te := &kasten.External{
				Sku: sku.External[kennung.Kasten, *kennung.Kasten]{
					Kennung: kennung.MustKasten(fdee.FileNameSansExt()),
				},
				FD: fdee,
			}

			//TODO-P2 offer option to edit again
			if _, err = formatText.Parse(f, &te.Objekte); err != nil {
				err = errors.Wrap(err)
				return
			}

			ts.Add(te)
		}()
	}

	out = ts.Copy()

	return
}

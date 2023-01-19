package commands

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/fd"
	"github.com/friedenberg/zit/src/foxtrot/id_set"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/mike/store_objekten"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/user_ops"
)

type EditEtikett struct {
	All bool
}

func init() {
	registerCommand(
		"edit-etikett",
		func(f *flag.FlagSet) Command {
			c := &EditEtikett{}

			f.BoolVar(&c.All, "all", false, "edit all Etiketten")

			return commandWithIds{CommandWithIds: c}
		},
	)
}

func (c EditEtikett) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet(
		id_set.ProtoId{
			Setter: &kennung.Etikett{},
		},
	)

	return
}

func (c EditEtikett) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	tks := ids.Etiketten.Copy()

	switch {
	case tks.Len() == 0 && !c.All:
		err = errors.Normalf("No Etiketten specified for editing. To edit all, use -all.")
		return

	case c.All && tks.Len() > 0:
		errors.Err().Print("Ignoring arguments because -all is set.")

		fallthrough

	case c.All:
		mtks := collections.MakeMutableValueSet[kennung.Etikett, *kennung.Etikett]()

		u.Konfig().Etiketten.Each(
			func(tt *etikett.Transacted) (err error) {
				return mtks.Add(tt.Sku.Kennung)
			},
		)

		tks = mtks.Copy()
	}

	var ps []string

	if ps, err = c.makeTempEtikettFiles(u, tks); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithCursorLocation(2, 3).
			WithFileType("zit-etikett").
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

	var tes collections.Set[*etikett.External]

	if tes, err = c.readTempEtikettFiles(u, ps); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	if err = tes.Each(
		func(e *etikett.External) (err error) {
			if _, err = u.StoreObjekten().Etikett().CreateOrUpdate(
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

func (c EditEtikett) makeTempEtikettFiles(
	u *umwelt.Umwelt,
	tks collections.ValueSet[kennung.Etikett, *kennung.Etikett],
) (ps []string, err error) {
	ps = make([]string, 0, tks.Len())

	var tempDir string

	if tempDir, err = files.TempDir(); err != nil {
		err = errors.Wrap(err)
		return
	}

	format := etikett.MakeFormatText(u.StoreObjekten())

	if err = tks.Each(
		func(tk kennung.Etikett) (err error) {
			var tt *etikett.Transacted

			if tt, err = u.StoreObjekten().Etikett().ReadOne(&tk); err != nil {
				if errors.Is(err, store_objekten.ErrNotFound{}) {
					err = nil
					tt = &etikett.Transacted{
						Sku: sku.Transacted[kennung.Etikett, *kennung.Etikett]{
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
				path.Join(tempDir, fmt.Sprintf("%s.%s", tk.String(), u.Konfig().FileExtensions.Etikett)),
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

func (c EditEtikett) readTempEtikettFiles(
	u *umwelt.Umwelt,
	ps []string,
) (out collections.Set[*etikett.External], err error) {
	ts := collections.MakeMutableSet[*etikett.External](
		etikett.ExternalKeyer{}.Key,
	)

	formatText := etikett.MakeFormatText(u.StoreObjekten())

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

			te := &etikett.External{
				Sku: sku.External[kennung.Etikett, *kennung.Etikett]{
					Kennung: kennung.MustEtikett(fdee.FileNameSansExt()),
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

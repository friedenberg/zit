package commands

import (
	"os"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CommandWithCwdQuery interface {
	RunWithCwdQuery(
		store *umwelt.Umwelt,
		ms kennung.MetaSet,
		cwdFiles *cwd.CwdFiles,
	) error
	DefaultGattungen() gattungen.Set
}

type commandWithCwdQuery struct {
	CommandWithCwdQuery
}

func (c commandWithCwdQuery) Complete(
	u *umwelt.Umwelt,
	args ...string,
) (err error) {
	var cgg CompletionGattungGetter
	ok := false

	if cgg, ok = c.CommandWithCwdQuery.(CompletionGattungGetter); !ok {
		return
	}

	cg := cgg.CompletionGattung()

	if cg.Contains(gattung.Zettel) {
		func() {
			zw := zettel.MakeWriterComplete(os.Stdout)
			defer errors.Deferred(&err, zw.Close)

			w := zw.WriteZettelVerzeichnisse

			if err = u.StoreObjekten().Zettel().ReadAllSchwanzen(w); err != nil {
				err = errors.Wrap(err)
				return
			}
		}()
	}

	if cg.Contains(gattung.Etikett) {
		if err = u.StoreObjekten().GetKennungIndex().EachSchwanzen(
			func(e kennung.IndexedLike[kennung.Etikett, *kennung.Etikett]) (err error) {
				if err = errors.Out().Printf("%s\tEtikett", e.GetKennung().String()); err != nil {
					err = errors.IsAsNilOrWrapf(
						err,
						syscall.EPIPE,
						"Etikett: %s",
						e.GetKennung(),
					)

					return
				}

				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if cg.Contains(gattung.Typ) {
		if err = u.Konfig().Typen.Each(
			func(tt typ.Transacted) (err error) {
				if err = errors.Out().Printf("%s\tTyp", tt.GetKennung()); err != nil {
					err = errors.IsAsNilOrWrapf(
						err,
						syscall.EPIPE,
						"Typ: %s",
						tt.GetKennung(),
					)

					return
				}

				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c commandWithCwdQuery) Run(u *umwelt.Umwelt, args ...string) (err error) {
	var cwdFiles cwd.CwdFiles

	if cwdFiles, err = cwd.MakeCwdFilesAll(
		u.Konfig(),
		u.Standort().Cwd(),
		u.StoreObjekten(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	ids := u.MakeMetaIdSetWithoutExcludedHidden(cwdFiles, c.DefaultGattungen())

	if err = ids.SetMany(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithCwdQuery(u, ids, &cwdFiles); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

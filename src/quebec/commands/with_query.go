package commands

import (
	"os"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/transacted"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CommandWithQuery interface {
	RunWithQuery(store *umwelt.Umwelt, ids kennung.MetaSet) error
	DefaultGattungen() gattungen.Set
}

type commandWithQuery struct {
	CommandWithQuery
}

type CompletionGattungGetter interface {
	CompletionGattung() gattungen.Set
}

func (c commandWithQuery) Complete(
	u *umwelt.Umwelt,
	args ...string,
) (err error) {
	var cgg CompletionGattungGetter
	ok := false

	if cgg, ok = c.CommandWithQuery.(CompletionGattungGetter); !ok {
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
			func(tt transacted.Typ) (err error) {
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

func (c commandWithQuery) Run(u *umwelt.Umwelt, args ...string) (err error) {
	ids := u.MakeMetaIdSetWithExcludedHidden(nil, c.DefaultGattungen())

	if err = ids.SetMany(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithQuery(u, ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

package commands

import (
	"os"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CommandWithQuery interface {
	RunWithQuery(store *umwelt.Umwelt, ids kennung.MetaSet) error
}

type commandWithQuery struct {
	CommandWithQuery
}

func (c commandWithQuery) Complete(u *umwelt.Umwelt, args ...string) (err error) {
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
		var ea []kennung.Etikett

		if ea, err = u.StoreObjekten().GetKennungIndex().GetAllEtiketten(); err != nil {
			err = errors.Wrap(err)
			return
		}

		for _, e := range ea {
			if err = errors.Out().Printf("%s\tEtikett", e.String()); err != nil {
				err = errors.IsAsNilOrWrapf(
					err,
					syscall.EPIPE,
					"Etikett: %s",
					e,
				)

				return
			}
		}

		return
	}

	if cg.Contains(gattung.Typ) {
		if err = u.Konfig().Typen.Each(
			func(tt *typ.Transacted) (err error) {
				if err = errors.Out().Printf("%s\tTyp", tt.Sku.Kennung); err != nil {
					err = errors.IsAsNilOrWrapf(
						err,
						syscall.EPIPE,
						"Typ: %s",
						tt.Sku.Kennung,
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
	ids := u.MakeMetaIdSet()

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
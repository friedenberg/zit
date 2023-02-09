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

type CommandWithHinweis interface {
	RunWithHinweis(store *umwelt.Umwelt, h kennung.Hinweis) error
}

type commandWithHinweis struct {
	CommandWithHinweis
}

func (c commandWithHinweis) Complete(u *umwelt.Umwelt, args ...string) (err error) {
	errors.TodoP0("implement")

	var cgg CompletionGattungGetter
	ok := false

	if cgg, ok = c.CommandWithHinweis.(CompletionGattungGetter); !ok {
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

func (c commandWithHinweis) Run(u *umwelt.Umwelt, args ...string) (err error) {
	errors.TodoP1("add metaid type to kennung and support for sigils")
	if len(args) != 0 {
		err = errors.Normalf("only one hinweis is accepted")
		return
	}

	var h kennung.Hinweis
	v := args[0]

	if h, err = u.StoreObjekten().GetAbbrStore().ExpandHinweisString(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithHinweis(u, h); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

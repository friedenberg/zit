package commands

import (
	"os"
	"syscall"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/charlie/gattung"
	"code.linenisgreat.com/zit-go/src/delta/gattungen"
	"code.linenisgreat.com/zit-go/src/echo/kennung"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
	"code.linenisgreat.com/zit-go/src/india/matcher"
	to_merge "code.linenisgreat.com/zit-go/src/india/sku_fmt"
	"code.linenisgreat.com/zit-go/src/kilo/cwd"
	"code.linenisgreat.com/zit-go/src/oscar/umwelt"
)

type CommandWithCwdQuery interface {
	RunWithCwdQuery(
		store *umwelt.Umwelt,
		ms matcher.Query,
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
			zw := to_merge.MakeWriterComplete(os.Stdout)
			defer errors.Deferred(&err, zw.Close)

			w := zw.WriteZettelVerzeichnisse

			if err = u.StoreObjekten().ReadAllSchwanzen(
				gattungen.MakeSet(gattung.Zettel),
				w,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}()
	}

	if cg.Contains(gattung.Etikett) {
		if err = u.StoreObjekten().GetKennungIndex().EachSchwanzen(
			func(e *kennung.IndexedLike) (err error) {
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
			func(tt *sku.Transacted) (err error) {
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
	ids := u.MakeMetaIdSetWithoutExcludedHidden(
		c.DefaultGattungen(),
	)

	if err = ids.SetMany(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithCwdQuery(u, ids, u.StoreUtil().GetCwdFiles()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

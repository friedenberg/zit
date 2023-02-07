package commands

import (
	"os"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CommandWithIds interface {
	RunWithIds(store *umwelt.Umwelt, ids kennung.Set) error
}

type CommandWithIdsAndProtoSet interface {
	CommandWithIds
	ProtoIdSet(*umwelt.Umwelt) kennung.ProtoIdSet
}

type commandWithIds struct {
	CommandWithIds
}

func (c commandWithIds) getIdProtoSet(u *umwelt.Umwelt) (is kennung.ProtoIdSet) {
	tid, hasCustomProtoSet := c.CommandWithIds.(CommandWithIdsAndProtoSet)

	switch {
	case hasCustomProtoSet:
		is = tid.ProtoIdSet(u)

	default:
		is = kennung.MakeProtoIdSet(
			kennung.ProtoId{
				Setter: &sha.Sha{},
			},
			kennung.ProtoId{
				Setter: &kennung.Hinweis{},
				Expand: func(v string) (out string, err error) {
					var h kennung.Hinweis
					h, err = u.StoreObjekten().GetAbbrStore().ExpandHinweisString(v)
					out = h.String()
					return
				},
			},
			kennung.ProtoId{
				Setter: &kennung.Typ{},
			},
			kennung.ProtoId{
				Setter: &ts.Time{},
			},
		)
	}

	return
}

func (c commandWithIds) Complete(u *umwelt.Umwelt, args ...string) (err error) {
	errors.TodoP0("implement")
	ps := c.getIdProtoSet(u)

	if ps.Contains(&kennung.Hinweis{}) {
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

	if ps.Contains(&kennung.Etikett{}) {
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

	if ps.Contains(&kennung.Typ{}) {
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

		return
	}

	return
}

func (c commandWithIds) Run(u *umwelt.Umwelt, args ...string) (err error) {
	ids := kennung.MakeSetWithExpanders(
		func(v string) (out string, err error) {
			var s sha.Sha
			s, err = u.StoreObjekten().GetAbbrStore().ExpandShaString(v)
			out = s.String()
			return
		},
		func(v string) (out string, err error) {
			var e kennung.Etikett
			e, err = u.StoreObjekten().GetAbbrStore().ExpandEtikettString(v)
			out = e.String()
			return
		},
		func(v string) (out string, err error) {
			var h kennung.Hinweis
			h, err = u.StoreObjekten().GetAbbrStore().ExpandHinweisString(v)
			out = h.String()
			return
		},
		nil, //typExpander func(string) (string, error),
		nil, //kastenExpander func(string) (string, error),
	)

	if err = ids.SetMany(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithIds(u, ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

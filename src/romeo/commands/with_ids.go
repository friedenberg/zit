package commands

import (
	"os"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/papa/umwelt"
)

type CommandWithIds interface {
	RunWithIds(store *umwelt.Umwelt, ids id_set.Set) error
}

type CommandWithIdsAndProtoSet interface {
	CommandWithIds
	ProtoIdSet(*umwelt.Umwelt) id_set.ProtoIdSet
}

type commandWithIds struct {
	CommandWithIds
}

func (c commandWithIds) getIdProtoSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	tid, hasCustomProtoSet := c.CommandWithIds.(CommandWithIdsAndProtoSet)

	switch {
	case hasCustomProtoSet:
		is = tid.ProtoIdSet(u)

	default:
		is = id_set.MakeProtoIdSet(
			id_set.ProtoId{
				MutableId: &sha.Sha{},
			},
			id_set.ProtoId{
				MutableId: &hinweis.Hinweis{},
				Expand: func(v string) (out string, err error) {
					var h hinweis.Hinweis
					h, err = u.StoreObjekten().Abbr().ExpandHinweisString(v)
					out = h.String()
					return
				},
			},
			id_set.ProtoId{
				MutableId: &kennung.Typ{},
			},
			id_set.ProtoId{
				MutableId: &ts.Time{},
			},
		)
	}

	return
}

func (c commandWithIds) Complete(u *umwelt.Umwelt, args ...string) (err error) {
	ps := c.getIdProtoSet(u)

	if ps.Contains(&hinweis.Hinweis{}) {
		func() {
			zw := zettel.MakeWriterComplete(os.Stdout)
			defer zw.Close()

			w := zw.WriteZettelTransacted

			if err = u.StoreObjekten().Zettel().ReadAllSchwanzenTransacted(w); err != nil {
				err = errors.Wrap(err)
				return
			}
		}()
	}

	if ps.Contains(&kennung.Etikett{}) {
		var ea []kennung.Etikett

		if ea, err = u.StoreObjekten().Etiketten(); err != nil {
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

	return
}

func (c commandWithIds) Run(u *umwelt.Umwelt, args ...string) (err error) {
	ps := c.getIdProtoSet(u)

	var ids id_set.Set

	if ids, err = ps.Make(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithIds(u, ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

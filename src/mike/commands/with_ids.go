package commands

import (
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/charlie/typ"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/kilo/umwelt"
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
					h, err = u.StoreObjekten().ExpandHinweisString(v)
					out = h.String()
					return
				},
			},
			id_set.ProtoId{
				MutableId: &typ.Typ{},
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
		var zts zettel_transacted.Set

		if zts, err = u.StoreObjekten().ZettelenSchwanzen(); err != nil {
			err = errors.Wrap(err)
			return
		}

		err = zts.Each(
			func(zt zettel_transacted.Zettel) (err error) {
        errors.PrintOutf(
          "%s\tZettel: !%s %s",
          zt.Named.Hinweis,
          zt.Named.Stored.Zettel.Typ,
          zt.Named.Stored.Zettel.Bezeichnung,
        )

				return
			},
		)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if ps.Contains(&etikett.Etikett{}) {
		var ea []etikett.Etikett

		if ea, err = u.StoreObjekten().Etiketten(); err != nil {
			err = errors.Wrap(err)
			return
		}

		for _, e := range ea {
			if err = errors.PrintOutf("%s\tEtikett", e.String()); err != nil {
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

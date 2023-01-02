package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/user_ops"
	"github.com/friedenberg/zit/src/remote_messages"
)

type Pull struct {
	gattung.Gattung
	All bool
}

func init() {
	registerCommand(
		"pull",
		func(f *flag.FlagSet) Command {
			c := &Pull{
				Gattung: gattung.Zettel,
			}

			f.Var(&c.Gattung, "gattung", "Gattung")
			f.BoolVar(&c.All, "all", false, "pull all Objekten")

			return c
		},
	)
}

func (c Pull) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	switch c.Gattung {

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
				MutableId: &kennung.Etikett{},
				Expand: func(v string) (out string, err error) {
					var e kennung.Etikett
					e, err = u.StoreObjekten().Abbr().ExpandEtikettString(v)
					out = e.String()
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

	case gattung.Typ:
		is = id_set.MakeProtoIdSet(
			id_set.ProtoId{
				MutableId: &kennung.Typ{},
			},
		)

	case gattung.Transaktion:
		is = id_set.MakeProtoIdSet(
			id_set.ProtoId{
				MutableId: &ts.Time{},
			},
		)
	}

	return
}

func (c Pull) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) == 0 {
		err = errors.Normalf("must specify kasten to pull from")
		return
	}

	from := args[0]

	if len(args) > 1 {
		args = args[1:]

		if c.All {
			errors.Log().Print("-all is set but arguments passed in. Ignore -all.")
		}
	} else if !c.All {
		err = errors.Normalf("Refusing to pull all unless -all is set.")
		return
	} else {
		args = []string{}
	}

	errors.Log().Print(args)

	ps := c.ProtoIdSet(u)

	var ids id_set.Set

	if ids, err = ps.Make(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	filter := id_set.Filter{
		AllowEmpty: c.All,
		Set:        ids,
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	var s *remote_messages.StageCommander

	if s, err = remote_messages.MakeStageCommander(
		u,
		from,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("")

	defer errors.Deferred(&err, s.Close)

	var dialoguePull remote_messages.Dialogue

	errors.Log().Printf("starting pull dialogue")

	if dialoguePull, err = s.StartDialogue(
		remote_messages.DialogueTypePush,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	pullOrPushOp := user_ops.MakeRemoteMessagesPullOrPush(u)

	if err = c.handleDialoguePull(
		pullOrPushOp,
		s,
		filter,
		dialoguePull,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Pull) handleDialoguePull(
	op user_ops.RemoteMessagesPullOrPush,
	s *remote_messages.StageCommander,
	filter id_set.Filter,
	d remote_messages.Dialogue,
) (err error) {
	if err = d.Send(filter); err != nil {
		err = errors.Wrap(err)
		return
	}

	t := transaktion.MakeTransaktion(ts.Now())

	if err = d.Receive(&t); err != nil {
		err = errors.Wrap(err)
		return
	}

	var pullObjektenDialogue remote_messages.Dialogue

	if pullObjektenDialogue, err = s.StartDialogue(
		remote_messages.DialogueTypePushObjekten,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	go errors.Deferred(
		&err,
		func() error {
			return op.HandleDialoguePullObjekten(s, pullObjektenDialogue, t.Skus)
		},
	)

	return
}

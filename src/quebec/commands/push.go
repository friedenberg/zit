package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Push struct {
	gattung.Gattung
	All bool
}

func init() {
	registerCommand(
		"push",
		func(f *flag.FlagSet) Command {
			c := &Push{
				Gattung: gattung.Zettel,
			}

			f.Var(&c.Gattung, "gattung", "Gattung")
			f.BoolVar(&c.All, "all", false, "pull all Objekten")

			return c
		},
	)
}

func (c Push) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
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

func (c Push) Run(u *umwelt.Umwelt, args ...string) (err error) {
	// if len(args) == 0 {
	// 	err = errors.Normalf("must specify kasten to pull from")
	// 	return
	// }

	// from := args[0]

	// if len(args) > 1 {
	// 	args = args[1:]

	// 	if c.All {
	// 		errors.Log().Print("-all is set but arguments passed in. Ignore -all.")
	// 	}
	// } else if !c.All {
	// 	err = errors.Normalf("Refusing to pull all unless -all is set.")
	// 	return
	// } else {
	// 	args = []string{}
	// }

	// errors.Log().Print(args)

	// ps := c.ProtoIdSet(u)

	// var ids id_set.Set

	// if ids, err = ps.Make(args...); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// filter := id_set.Filter{
	// 	AllowEmpty: c.All,
	// 	Set:        ids,
	// }

	// if err = u.Lock(); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// defer errors.Deferred(&err, u.Unlock)

	// var s *remote_messages.StageCommander

	// if s, err = remote_messages.MakeStageCommander(
	// 	u,
	// 	from,
	// ); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// errors.Log().Printf("")

	// defer errors.Deferred(&err, s.Close)

	// var dialoguePush remote_messages.Dialogue

	// errors.Log().Printf("starting pull dialogue")

	// if dialoguePush, err = s.StartDialogue(
	// 	remote_messages.DialogueTypePull,
	// ); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// pullOrPushOp := user_ops.MakeRemoteMessagesPullOrPush(u)

	// if err = pullOrPushOp.HandleDialoguePush(
	// 	s,
	// 	filter,
	// 	dialoguePush,
	// ); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	return
}

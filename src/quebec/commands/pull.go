package commands

import (
	"flag"
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/oscar/umwelt"
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
		//TODO-P3 handle all is set
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
		remote_messages.DialogueTypePull,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.handleDialoguePull(
		u,
		s,
		filter,
		dialoguePull,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.MainDialogue().Send(remote_messages.MessageDone{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Pull) handleDialoguePullObjekten(
	u *umwelt.Umwelt,
	s *remote_messages.StageCommander,
	d remote_messages.Dialogue,
	skus sku.MutableSet,
	wg *sync.WaitGroup,
) (err error) {
	skusNeeded := sku.MakeMutableSet()

	if err = skus.Each(
		func(sk *sku.Sku) (err error) {
			//TODO-P1 support other gattung
			if sk.Gattung != gattung.Zettel {
				return
			}

			if u.StoreObjekten().Zettel().HasObjekte(sk.Sha) {
				return
			}

			skusNeeded.Add(*sk)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer wg.Done()

	if err = d.Send(skusNeeded); err != nil {
		errors.Log().Print("failed to send skus")
		err = errors.Wrap(err)
		return
	}

	errors.Log().Print("did send skus")

	for {
		errors.Log().Print("did start loop")
		var zt zettel.Transacted

		if err = d.Receive(&zt); err != nil {
			errors.Log().Printf("did receive error: %s", errors.Wrap(err))

			if errors.IsEOF(err) {
				err = nil
				break
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := c.handleDialoguePullAkte(u, s, zt.Objekte.Akte); err != nil {
				errors.Log().Printf("pull akte error: %s", err)
			}
		}()

		if err = u.StoreObjekten().Zettel().Inherit(&zt); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.Log().Print("did receive no error")

		errors.Log().Print(zt)
	}

	return
}

func (c Pull) handleDialoguePullAkte(
	u *umwelt.Umwelt,
	s *remote_messages.StageCommander,
	sh sha.Sha,
) (err error) {
	if sh.IsNull() {
		return
	}

	var d remote_messages.Dialogue

	if d, err = s.StartDialogue(remote_messages.DialogueTypePullAkte); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, d.Close)

	if err = d.Send(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	var aw sha.WriteCloser

	if aw, err = u.StoreObjekten().AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, aw.Close)

	if _, err = io.Copy(aw, d); err != nil {
		if errors.IsEOF(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Pull) handleDialoguePull(
	u *umwelt.Umwelt,
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
		remote_messages.DialogueTypePullObjekten,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go c.handleDialoguePullObjekten(u, s, pullObjektenDialogue, t.Skus, wg)
	wg.Wait()

	//	//TODO-P2 deal with errors that might close the channel
	//	if err = u.StoreObjekten().Zettel().Inherit(z); err != nil {
	//		err = errors.Wrap(err)
	//		return
	//	}
	//}

	return
}

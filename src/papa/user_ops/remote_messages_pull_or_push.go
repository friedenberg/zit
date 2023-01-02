package user_ops

import (
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/remote_messages"
)

type RemoteMessagesPullOrPush struct {
	pullObjektenWaitGroup *sync.WaitGroup
	umwelt                *umwelt.Umwelt
}

func MakeRemoteMessagesPullOrPush(
	u *umwelt.Umwelt,
) RemoteMessagesPullOrPush {
	return RemoteMessagesPullOrPush{
		pullObjektenWaitGroup: &sync.WaitGroup{},
		umwelt:                u,
	}
}

func (op RemoteMessagesPullOrPush) AddToSoldierStage(
	s *remote_messages.StageSoldier,
) {
	s.RegisterHandler(
		remote_messages.DialogueTypePush,
		op.handleDialoguePush,
	)

	s.RegisterHandler(
		remote_messages.DialogueTypePushObjekten,
		op.handleDialoguePushObjekten,
	)

	s.RegisterHandler(
		remote_messages.DialogueTypePushAkte,
		op.handleDialoguePushAkte,
	)

	// s.RegisterHandler(
	// 	remote_messages.DialogueTypePull,
	// 	op.HandleDialoguePull,
	// )

	// s.RegisterHandler(
	// 	remote_messages.DialogueTypePullObjekten,
	// 	op.handleDialoguePullObjekten,
	// )

	// s.RegisterHandler(
	// 	remote_messages.DialogueTypePullAkte,
	// 	op.handleDialoguePullAkte,
	// )
}

func (op RemoteMessagesPullOrPush) handleDialoguePushAkte(
	d remote_messages.Dialogue,
) (err error) {
	var sh sha.Sha

	errors.Log().Print("waiting to receive sha")

	if err = d.Receive(&sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, d.Close)

	errors.Log().Printf("did receive sha: %s", sh)

	var ar sha.ReadCloser

	if ar, err = op.umwelt.StoreObjekten().AkteReader(sh); err != nil {
		errors.Log().Printf("got error on akte reader: %s", err)
		if errors.IsNotExist(err) {
			err = nil
			return
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	var n int64

	if n, err = io.Copy(d, ar); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("copied %d bytes", n)

	return
}

func (op RemoteMessagesPullOrPush) SendFiltered(
	d remote_messages.Dialogue,
	filter id_set.Filter,
) (err error) {
	t := transaktion.MakeTransaktion(ts.Now())

	method := op.umwelt.StoreObjekten().Zettel().ReadAllSchwanzenVerzeichnisse

	if op.umwelt.Konfig().IncludeHistory {
		method = op.umwelt.StoreObjekten().Zettel().ReadAllVerzeichnisse
	}

	if err = method(
		collections.MakeChain(
			zettel.WriterIds{Filter: filter}.WriteZettelVerzeichnisse,
			func(z *zettel.Transacted) (err error) {
				t.Skus.Add2(&z.Sku)
				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.Send(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op RemoteMessagesPullOrPush) handleDialoguePush(
	d remote_messages.Dialogue,
) (err error) {
	var filter id_set.Filter

	if err = d.Receive(&filter); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = op.SendFiltered(d, filter); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op RemoteMessagesPullOrPush) SendSkus(
	d remote_messages.Dialogue,
	skus []sku.Sku,
) (err error) {
	errors.Log().Printf("starting zettel send loop: %d", len(skus))

	for _, s := range skus {
		//TODO-P1 support any transacted objekte
		if s.Gattung != gattung.Zettel {
			errors.Err().Printf("not a zettel, continuing: %v", s)
			continue
		}

		errors.Log().Printf("found a zettel, sending: %v", s)

		var zt *zettel.Transacted

		if zt, err = op.umwelt.StoreObjekten().Zettel().Inflate(
			//TODO-P2 use an actually correct time
			ts.Now(),
			&s,
		); err != nil {
			errors.Log().Printf("error inflating zettel for send: %v", err)
			err = errors.Wrap(err)
			return
		}

		errors.Log().Printf("sending zettel.Transacted: %v", zt)

		if err = d.Send(zt); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.Log().Printf("did send zettel.Transacted: %v", zt)
	}

	return
}

func (op RemoteMessagesPullOrPush) handleDialoguePushObjekten(
	d remote_messages.Dialogue,
) (err error) {
	var skus []sku.Sku

	errors.Log().Print("waiting to receive skus")

	if err = d.Receive(&skus); err != nil {
		errors.Log().Printf("error receiving skus: %s", err)
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("did receive skus: %d", len(skus))

	if err = op.SendSkus(d, skus); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("send objekten")

	if err = d.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("did close")

	return
}

func (op RemoteMessagesPullOrPush) SendNeededSkus(
	d remote_messages.Dialogue,
	skus sku.MutableSet,
) (err error) {
	skusNeeded := collections.MakeMutableSet[*sku.Sku](
		func(sk *sku.Sku) string {
			if sk == nil {
				return ""
			}

			return sk.Sha.String()
		},
	)

	if err = skus.Each(
		func(sk *sku.Sku) (err error) {
			//TODO-P1 support other gattung
			if sk.Gattung != gattung.Zettel {
				return
			}

			if op.umwelt.StoreObjekten().Zettel().HasObjekte(sk.Sha) {
				errors.Log().Printf("already have %s", sk.Sha)
				return
			}

			errors.Log().Printf("don't have %s", sk.Sha)
			skusNeeded.Add(sk)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.Send(skusNeeded.Elements()); err != nil {
		errors.Log().Print("failed to send skus")
		err = errors.Wrap(err)
		return
	}

	return
}

func (op RemoteMessagesPullOrPush) HandleDialoguePullObjekten(
	s *remote_messages.StageCommander,
	d remote_messages.Dialogue,
	skus sku.MutableSet,
) (err error) {
	errors.Log().Print("did send skus")

	if err = op.SendNeededSkus(d, skus); err != nil {
		err = errors.Wrap(err)
		return
	}

	for {
		errors.Log().Print("did start zettel receive loop")
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

		go func() {
			if err := op.handleDialoguePullAkte(s, zt.Objekte.Akte); err != nil {
				errors.Log().Printf("pull akte error: %s", err)
			}
		}()

		if err = op.umwelt.StoreObjekten().Zettel().Inherit(&zt); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.Log().Print("did receive no error")
		errors.Log().Print(zt)
	}

	return
}

func (op RemoteMessagesPullOrPush) ReceiveAkte(
	d remote_messages.Dialogue,
) (err error) {
	var aw sha.WriteCloser

	if aw, err = op.umwelt.StoreObjekten().AkteWriter(); err != nil {
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

func (op RemoteMessagesPullOrPush) handleDialoguePullAkte(
	s *remote_messages.StageCommander,
	sh sha.Sha,
) (err error) {
	if sh.IsNull() {
		return
	}

	var d remote_messages.Dialogue

	if d, err = s.StartDialogue(remote_messages.DialogueTypePushAkte); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, d.Close)

	if err = d.Send(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = op.ReceiveAkte(d); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

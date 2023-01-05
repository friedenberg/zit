package user_ops

import (
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/id"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/remote_conn"
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
	s *remote_conn.StageSoldier,
) {
	s.RegisterHandler(
		remote_conn.DialogueTypeSkusForFilter,
		op.skusForFilter,
	)

	s.RegisterHandler(
		remote_conn.DialogueTypeObjekteReaderForSku,
		op.objekteReaderForSku,
	)

	s.RegisterHandler(
		remote_conn.DialogueTypePush,
		op.handleDialoguePush,
	)

	s.RegisterHandler(
		remote_conn.DialogueTypePushObjekten,
		op.handleDialoguePushObjekten,
	)

	s.RegisterHandler(
		remote_conn.DialogueTypePushAkte,
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
	d remote_conn.Dialogue,
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

	errors.Log().Printf("sent %d akte bytes", n)

	return
}

func (op RemoteMessagesPullOrPush) SendFiltered(
	d remote_conn.Dialogue,
	filter id_set.Filter,
) (err error) {
	method := op.umwelt.StoreObjekten().Zettel().ReadAllSchwanzenVerzeichnisse

	if op.umwelt.Konfig().IncludeHistory {
		method = op.umwelt.StoreObjekten().Zettel().ReadAllVerzeichnisse
	}

	if err = method(
		collections.MakeChain(
			zettel.WriterIds{Filter: filter}.WriteZettelVerzeichnisse,
			func(z *zettel.Transacted) (err error) {
				if err = d.Send(sku.String(&z.Sku)); err != nil {
					err = errors.Wrap(err)
					return
				}

				errors.Err().Print(z.Sku)

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op RemoteMessagesPullOrPush) objekteReaderForSku(
	d remote_conn.Dialogue,
) (err error) {
	var strSku string

	if err = d.Receive(&strSku); err != nil {
		err = errors.Wrap(err)
		return
	}

	var sk sku.SkuLike

	if sk, err = sku.MakeSku(strSku); err != nil {
		err = errors.Wrap(err)
		return
	}

	var or io.ReadCloser

	if or, err = op.umwelt.StoreObjekten().ReadCloserObjekten(
		id.Path(sk.GetObjekteSha(), op.umwelt.Standort().DirObjektenZettelen()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = io.Copy(d, or); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op RemoteMessagesPullOrPush) skusForFilter(
	d remote_conn.Dialogue,
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

func (op RemoteMessagesPullOrPush) handleDialoguePush(
	d remote_conn.Dialogue,
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

func (op RemoteMessagesPullOrPush) SendObjekte(
	d remote_conn.Dialogue,
	s sku.SkuLike,
) (err error) {
	//TODO-P1 support any transacted objekte
	if s.GetGattung() != gattung.Zettel {
		errors.Err().Printf("not a zettel, continuing: %v", s)
		return
	}

	sk := s.(*sku.Transacted[hinweis.Hinweis, *hinweis.Hinweis])

	errors.Log().Printf("found a zettel, sending: %v", s)

	var zt *zettel.Transacted

	if zt, err = op.umwelt.StoreObjekten().Zettel().Inflate(
		//TODO-P2 use an actually correct time
		ts.Now(),
		sk,
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

	return
}

func (op RemoteMessagesPullOrPush) handleDialoguePushObjekten(
	d remote_conn.Dialogue,
) (err error) {
	errors.Log().Print("waiting to receive skus")

	for {
		var strSku string

		if err = d.Receive(&strSku); err != nil {
			errors.Log().Printf("error receiving skus: %s", err)
			err = errors.Wrap(err)
			return
		}

		errors.Log().Printf("received sku: %s", strSku)

		var sk sku.SkuLike

		if sk, err = sku.MakeSku(strSku); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.Log().Printf("sending sku: %s", strSku)

		if err = op.SendObjekte(d, sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.Log().Printf("did send sku: %s", strSku)
	}
}

func (op RemoteMessagesPullOrPush) SendNeededSkus(
	d remote_conn.Dialogue,
	skus []sku.SkuLike,
) (err error) {
	errors.Log().Print("sending needed skus")
	defer errors.Log().Print("sent needed skus")

	for _, sk := range skus {
		if err = d.Send(sku.String(sk)); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (op RemoteMessagesPullOrPush) HandleDialoguePullObjekten(
	s *remote_conn.StageCommander,
	d remote_conn.Dialogue,
	skus []sku.SkuLike,
) (err error) {
	if err = op.SendNeededSkus(d, skus); err != nil {
		err = errors.Wrap(err)
		return
	}

	for i := 0; i < len(skus); i++ {
		errors.Log().Print("waiting to receive one zettel")
		var zt zettel.Transacted

		if err = d.Receive(&zt); err != nil {
			errors.Log().Printf("did receive error: %s", errors.Wrap(err))

			err = errors.Wrap(err)
			return
		}

		errors.Log().Printf("did receive one zettel: %s", zt.Sku.Kennung)

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
	d remote_conn.Dialogue,
) (err error) {
	var aw sha.WriteCloser

	if aw, err = op.umwelt.StoreObjekten().AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, aw.Close)

	var n int64

	if n, err = io.Copy(aw, d); err != nil {
		if errors.IsEOF(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	errors.Log().Printf("received %d akte bytes", n)

	return
}

func (op RemoteMessagesPullOrPush) handleDialoguePullAkte(
	s *remote_conn.StageCommander,
	sh sha.Sha,
) (err error) {
	if sh.IsNull() {
		return
	}

	var d remote_conn.Dialogue

	if d, err = s.StartDialogue(remote_conn.DialogueTypePushAkte); err != nil {
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

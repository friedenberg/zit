package remote_push

import (
	"syscall"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/juliett/query"
	"code.linenisgreat.com/zit/src/november/umwelt"
	"code.linenisgreat.com/zit/src/papa/remote_conn"
)

const (
	concurrentSkuFilterJobLimit = 100
	// concurrentSkuFilterJobLimit = 1
)

type Client interface {
	SendNeededSkus(*query.Group) error
	Close() error
}

type client struct {
	umwelt             *umwelt.Umwelt
	stage              *remote_conn.StageCommander
	chDone             chan struct{}
	chFilterSkuTickets chan struct{}
}

func MakeClient(u *umwelt.Umwelt, from string) (c *client, err error) {
	c = &client{
		chDone:             make(chan struct{}),
		chFilterSkuTickets: make(chan struct{}, concurrentSkuFilterJobLimit),
	}

	if c.stage, err = remote_conn.MakeStageCommander(
		u,
		from,
		"push",
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c client) Close() (err error) {
	if err = c.stage.MainDialogue().Send(struct{}{}); err != nil {
		if errors.IsErrno(err, syscall.EPIPE) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if err = c.stage.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c client) SendNeededSkus(filter *query.Group) (err error) {
	err = todo.Implement()
	// var d remote_conn.Dialogue

	// if d, err = c.stage.StartDialogue(
	// 	remote_conn.DialogueTypeGetNeededSkus,
	// ); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// if err = d.Send(in); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// if err = d.Receive(out); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	return
}

func (c *client) ObjekteWriter(
	g schnittstellen.GattungGetter,
	sh schnittstellen.ShaLike,
) (rc sha.ReadCloser, err error) {
	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeObjekteWriter,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	msgRequest := messageRequestObjekteData{
		Gattung: gattung.Make(g.GetGattung()),
	}

	msgRequest.Sha.SetShaLike(sh)

	if err = d.Send(msgRequest); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P1 copy objekte data

	return
}

func (c client) AkteWriter(
	sh sha.Sha,
) (rc sha.ReadCloser, err error) {
	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeAkteWriter,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.Send(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P1 copy akte data

	return
}

package remote_push

import (
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/remote_conn"
	"github.com/friedenberg/zit/src/schnittstellen"
)

const (
	concurrentSkuFilterJobLimit = 100
	// concurrentSkuFilterJobLimit = 1
)

type Client interface {
	SendNeededSkus(id_set.Filter) error
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

// TODO-P0
func (c client) SendNeededSkus(filter id_set.Filter) (err error) {
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
	sh gattung.ShaLike,
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
		Sha:     sha.Make(sh.GetSha()),
	}

	if err = d.Send(msgRequest); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO-P0 copy objekte data

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

	//TODO-P0 copy akte data

	return
}

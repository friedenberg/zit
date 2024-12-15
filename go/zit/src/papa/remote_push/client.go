package remote_push

import (
	"syscall"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/remote_conn"
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
	env                *env.Local
	stage              *remote_conn.StageCommander
	chDone             chan struct{}
	chFilterSkuTickets chan struct{}
}

func MakeClient(u *env.Local, from string) (c *client, err error) {
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

func (c *client) ObjectWriter(
	g interfaces.GenreGetter,
	sh interfaces.Sha,
) (rc sha.ReadCloser, err error) {
	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeObjectWriter,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	msgRequest := messageRequestObjectData{
		Gattung: genres.Make(g.GetGenre()),
	}

	msgRequest.Sha.SetShaLike(sh)

	if err = d.Send(msgRequest); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P1 copy object data

	return
}

func (c client) BlobWriter(
	sh sha.Sha,
) (rc sha.ReadCloser, err error) {
	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeBlobWriter,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.Send(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P1 copy blob data

	return
}

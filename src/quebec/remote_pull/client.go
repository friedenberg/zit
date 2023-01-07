package remote_pull

import (
	"net"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/remote_conn"
)

type FuncSku func(sku.SkuLike) error

type Client interface {
	SkusFromFilter(id_set.Filter, FuncSku) error
	ObjekteReaderForSku(sku.SkuLike) (sha.ReadCloser, error)
	AkteReader(sha.Sha) (sha.ReadCloser, error)
	Close() error
}

type client struct {
	umwelt *umwelt.Umwelt
	stage  *remote_conn.StageCommander
	wg     *sync.WaitGroup
	chDone chan struct{}
}

func MakeClient(u *umwelt.Umwelt, from string) (c *client, err error) {
	c = &client{
		wg:     &sync.WaitGroup{},
		chDone: make(chan struct{}),
	}

	if c.stage, err = remote_conn.MakeStageCommander(
		u,
		from,
		"pull",
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c client) Close() (err error) {
	c.wg.Wait()

	if err = c.stage.MainDialogue().Send(struct{}{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.stage.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c client) SkusFromFilter(ids id_set.Filter, f FuncSku) (err error) {
	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeSkusForFilter,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.Send(ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	chErr := make(chan error)

	go func() {
		errors.DeferredChanError(&err, chErr)
		close(c.chDone)
		d.Close()
	}()

LOOP:
	for {
		select {
		case <-c.chDone:
			break LOOP

		default:
			break
		}

		var strSku string

		var err1 error

		if err1 = d.Receive(&strSku); err1 != nil {
			if errors.IsEOF(err1) || errors.Is(err1, net.ErrClosed) {
				err1 = nil
				break LOOP
			} else {
				err1 = errors.Wrap(err)
				chErr <- err1
				return
			}
		}

		errors.Log().Printf("received sku: %s", strSku)

		c.wg.Add(1)
		go c.makeAndProcessOneSkuStringWithFilter(
			strSku,
			f,
			chErr,
		)
	}

	return
}

var count int
var lock sync.Locker

func init() {
	lock = &sync.Mutex{}
}

func (c *client) makeAndProcessOneSkuStringWithFilter(
	strSku string,
	f FuncSku,
	chError chan<- error,
) {
	defer c.wg.Done()

	var err error
	var sk sku.SkuLike

	if sk, err = sku.MakeSku(strSku); err != nil {
		select {
		case <-c.chDone:
		case chError <- errors.Wrap(err):
		default:
			errors.Err().Printf("Client Error: %s", err)
		}

		return
	}

	if err = f(sk); err != nil {
		if errors.IsEOF(err) || errors.Is(err, net.ErrClosed) {
		} else {
			select {
			case <-c.chDone:
			case chError <- errors.Wrap(err):
			default:
				errors.Err().Printf("Client Error: %s", err)
			}
		}

		return
	}
}

func (c *client) ObjekteReaderForSku(
	sk sku.SkuLike,
) (rc sha.ReadCloser, err error) {
	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeObjekten,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.Send(sku.String(sk)); err != nil {
		err = errors.Wrap(err)
		return
	}

	rc = sha.MakeReadCloser(d)

	return
}

func (c client) AkteReader(
	sh sha.Sha,
) (rc sha.ReadCloser, err error) {
	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeAkten,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.Send(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	rc = sha.MakeReadCloser(d)

	return
}

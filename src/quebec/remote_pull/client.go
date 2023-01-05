package remote_pull

import (
	"bufio"
	"net"
	"strconv"
	"strings"
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

	chObjekteReaders  chan sha.ReadCloser
	lockObjekten1     sync.Locker
	lockObjekten2     sync.Locker
	dialogueObjekten  remote_conn.Dialogue
	onceObjekteReader *sync.Once

	onceAkteReader *sync.Once
}

func MakeClient(u *umwelt.Umwelt, from string) (c *client, err error) {
	c = &client{
		wg:                &sync.WaitGroup{},
		chDone:            make(chan struct{}),
		chObjekteReaders:  make(chan sha.ReadCloser),
		lockObjekten1:     &sync.Mutex{},
		lockObjekten2:     &sync.Mutex{},
		onceObjekteReader: &sync.Once{},
		onceAkteReader:    &sync.Once{},
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

	if err = c.stage.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c client) acquireConnectionLicense() {
	c.wg.Add(1)
	remote_conn.AcquireConnLicense()
}

func (c client) releaseConnectionLicense() {
	c.wg.Done()
	remote_conn.ReleaseConnLicense()
}

func (c client) SkusFromFilter(ids id_set.Filter, f FuncSku) (err error) {
	c.acquireConnectionLicense()
	defer c.releaseConnectionLicense()

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
				break
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

func (c client) makeAndProcessOneSkuStringWithFilter(
	strSku string,
	f FuncSku,
	chError chan<- error,
) {
	defer c.wg.Done()

	var err error
	var sk sku.SkuLike

	if sk, err = sku.MakeSku(strSku); err != nil {
		chError <- errors.Wrap(err)
		return
	}

	if err = f(sk); err != nil {
		chError <- errors.Wrap(err)
		return
	}
}

func (c *client) requestObjekten() {
	var err error

	if c.dialogueObjekten, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeObjekten,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	br := bufio.NewReader(c.dialogueObjekten)

	go func() {
		lContinue := &sync.Mutex{}

	LOOP:
		for {
			select {
			case <-c.chDone:
				break LOOP

			default:
				break
			}

			lContinue.Lock()

			var line string

			if line, err = br.ReadString('\n'); err != nil {
				if errors.IsEOF(err) || errors.Is(err, net.ErrClosed) {
					err = nil
					break
				} else {
					panic(err)
				}
			}

			var n int64

			if n, err = strconv.ParseInt(strings.TrimSpace(line), 10, 64); err != nil {
				errors.Log().Printf("%s", err)
				panic(err)
			}

			errors.Log().Printf("about to receive %d bytes", n)

			c.chObjekteReaders <- makeBoundReader(
				br,
				lContinue,
				n,
			)
		}
	}()
}

func (c *client) ObjekteReaderForSku(
	sk sku.SkuLike,
) (rc sha.ReadCloser, err error) {
	c.onceObjekteReader.Do(c.requestObjekten)

	c.lockObjekten1.Lock()

	if err = c.dialogueObjekten.Send(sku.String(sk)); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("send sku: %s", sk)

	c.lockObjekten2.Lock()
	c.lockObjekten1.Unlock()

	errors.Log().Printf("return reader")

	rc = <-c.chObjekteReaders

	c.lockObjekten2.Unlock()

	return
}

func (c client) AkteReader(
	sh sha.Sha,
) (rc sha.ReadCloser, err error) {
	c.acquireConnectionLicense()
	defer c.releaseConnectionLicense()

	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeAkteReaderForSha,
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

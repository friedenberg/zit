package commands

import (
	"flag"
	"io"
	"os"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type WriteBlob struct{}

func init() {
	registerCommand(
		"write-blob",
		func(f *flag.FlagSet) Command {
			c := &WriteBlob{}

			return c
		},
	)
}

type answer struct {
	sha.Sha
	Path string
}

func (c WriteBlob) Run(u *env.Env, args ...string) (err error) {
	wg := &sync.WaitGroup{}
	wg.Add(len(args))
	chCancel := make(chan struct{})
	chError := make(chan error)

	sawStdin := false

	for _, a := range args {
		switch {
		case sawStdin:
			ui.Err().Print("'-' passed in more than once. Ignoring")
			continue

		case a == "-":
			sawStdin = true
		}

		go c.doOne(
			chCancel,
			chError,
			wg,
			u.GetFSHome(),
			a,
		)
	}

	go func() {
		err = <-chError
		ui.Err().Print(err)
		close(chCancel)
	}()

	wg.Wait()

	return
}

func (c WriteBlob) doOne(
	chCancel <-chan struct{},
	chError chan<- error,
	wg *sync.WaitGroup,
	arf interfaces.BlobWriterFactory,
	p string,
) {
	var err error

	defer wg.Done()

	isDone := func() bool {
		select {
		case <-chCancel:
			return true

		default:
			return false
		}
	}

	var rc io.ReadCloser

	if p == "-" {
		rc = os.Stdin
	} else {
		if rc, err = files.Open(p); err != nil {
			err = errors.Wrap(err)
			chError <- err
			return
		}
	}

	if isDone() {
		return
	}

	defer errors.DeferredChan(chError, rc.Close)

	var wc fs_home.Writer

	if wc, err = arf.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		chError <- err
		return
	}

	defer errors.DeferredChan(chError, wc.Close)

	if isDone() {
		return
	}

	if _, err = io.Copy(wc, rc); err != nil {
		err = errors.Wrap(err)
		chError <- err
		return
	}

	if isDone() {
		return
	}

	ui.Out().Printf("%s %s", wc.GetShaLike(), p)
}

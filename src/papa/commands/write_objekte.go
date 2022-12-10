package commands

import (
	"flag"
	"io"
	"os"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/sha"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type WriteObjekte struct{}

func init() {
	registerCommand(
		"write-objekte",
		func(f *flag.FlagSet) Command {
			c := &WriteObjekte{}

			return c
		},
	)
}

type answer struct {
	sha.Sha
	Path string
}

func (c WriteObjekte) Run(u *umwelt.Umwelt, args ...string) (err error) {
	wg := &sync.WaitGroup{}
	wg.Add(len(args))
	chCancel := make(chan struct{})
	chError := make(chan error)

	sawStdin := false

	for _, a := range args {
		switch {
		case sawStdin:
			errors.Err().Print("'-' passed in more than once. Ignoring")
			continue

		case a == "-":
			sawStdin = true
		}

		go c.doOne(
			chCancel,
			chError,
			wg,
			u.StoreObjekten(),
			a,
		)
	}

	go func() {
		err = <-chError
		errors.Err().Print(err)
		close(chCancel)
	}()

	wg.Wait()

	return
}

func (c WriteObjekte) doOne(
	chCancel <-chan struct{},
	chError chan<- error,
	wg *sync.WaitGroup,
	arf gattung.AkteWriterFactory,
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

	var wc age_io.Writer

	if wc, err = arf.AkteWriter(); err != nil {
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

	errors.Out().Printf("%s %s", wc.Sha(), p)
}

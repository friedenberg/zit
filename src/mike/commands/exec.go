package commands

import (
	"flag"
	"io"
	"os"
	"os/exec"
	"path"
	"sync"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/delta/konfig"
	"github.com/friedenberg/zit/src/echo/id_set"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
)

type Exec struct {
}

func init() {
	registerCommand(
		"exec",
		func(f *flag.FlagSet) Command {
			c := &Exec{}

			return c
		},
	)
}

func (c Exec) Run(u *umwelt.Umwelt, args ...string) (err error) {
	var tz zettel_transacted.Transacted
	var executor konfig.RemoteScript
	var ar io.ReadCloser

	if tz, ar, executor, err = c.getZettel(u, args[0]); err != nil {
		err = errors.Error(err)
		return
	}

	var pipePath string

	if pipePath, err = c.makeFifoPipe(tz); err != nil {
		err = errors.Error(err)
		return
	}

	var cmd *exec.Cmd

	if cmd, err = c.makeCmd(executor, pipePath, args...); err != nil {
		err = errors.Error(err)
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go c.feedPipe(ar, wg, pipePath, tz)
	go c.exec(wg, pipePath, cmd)

	wg.Wait()
	logz.Print("done waiting")

	return
}

func (c Exec) getZettel(
	u *umwelt.Umwelt,
	hString string,
) (
	tz zettel_transacted.Transacted,
	ar io.ReadCloser,
	executor konfig.RemoteScript,
	err error,
) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(u); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	ps := id_set.MakeProtoSet(
		&sha.Sha{},
		&hinweis.Hinweis{},
		&hinweis.HinweisWithIndex{},
	)

	is := ps.MakeOne(hString)

	var idd id.Id
	ok := false

	if idd, ok = is.AnyShaOrHinweis(); !ok {
		err = errors.Errorf("unsupported id: %s", is)
		return
	}

	if tz, err = store.StoreObjekten().Read(idd); err != nil {
		err = errors.Error(err)
		return
	}

	typ := tz.Named.Stored.Zettel.Typ.String()

	if typKonfig, ok := store.Umwelt.Konfig.Typen[typ]; ok {
		executor = typKonfig.ExecCommand
	} else {
		err = errors.Normal(errors.Errorf("Typ does not have an exec-command set: %s", typ))
		return
	}

	if ar, err = store.StoreObjekten().AkteReader(tz.Named.Stored.Zettel.Akte); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (c Exec) makeFifoPipe(tz zettel_transacted.Transacted) (p string, err error) {
	h := tz.Named.Hinweis
	var d string

	if d, err = os.MkdirTemp("", h.Kopf()); err != nil {
		err = errors.Error(err)
		return
	}

	p = path.Join(d, h.Schwanz()+"."+tz.Named.Stored.Zettel.Typ.String())

	if err = syscall.Mknod(p, syscall.S_IFIFO|0666, 0); err != nil {
		err = errors.Error(err)
		return
	}

	if err = os.Remove(p); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (c Exec) makeCmd(
	executor konfig.RemoteScript,
	p string,
	args ...string,
) (cmd *exec.Cmd, err error) {
	cmdArgs := append([]string{p})

	if len(args) > 1 {
		cmdArgs = append(cmdArgs, args[1:]...)
	}

	if cmd, err = executor.Cmd(cmdArgs...); err != nil {
		err = errors.Error(err)
		return
	}

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return
}

func (c Exec) feedPipe(
	ar io.ReadCloser,
	wg *sync.WaitGroup,
	p string,
	tz zettel_transacted.Transacted,
) (err error) {
	defer wg.Done()
	var pipeFileWriter *os.File

	if pipeFileWriter, err = os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777); err != nil {
		err = errors.Error(err)
		return
	}

	defer pipeFileWriter.Close()

	defer ar.Close()

	if _, err = io.Copy(pipeFileWriter, ar); err != nil {
		err = errors.Error(err)
		return
	}

	logz.Print("done copying")

	return
}

func (c Exec) exec(
	wg *sync.WaitGroup,
	p string,
	cmd *exec.Cmd,
) (err error) {
	defer wg.Done()
	// var pipeFileReader *os.File

	// if pipeFileReader, err = os.OpenFile(pipePath, os.O_CREATE, os.ModeNamedPipe); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	// defer pipeFileReader.Close()

	// cmd.ExtraFiles = append(cmd.ExtraFiles, pipeFileReader)

	logz.Print("start running")

	if err = cmd.Run(); err != nil {
		err = errors.Error(err)
		return
	}

	logz.Print("done running")

	return
}

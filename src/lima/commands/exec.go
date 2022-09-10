package commands

import (
	"flag"
	"io"
	"os"
	"os/exec"
	"path"
	"sync"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/id"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/store_with_lock"
	"github.com/friedenberg/zit/src/juliett/umwelt"
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
	var tz zettel_transacted.Zettel
	var executor konfig.RemoteScript
	var ar io.ReadCloser

	if tz, ar, executor, err = c.getZettel(u, args[0]); err != nil {
		err = errors.Wrap(err)
		return
	}

	var pipePath string

	if pipePath, err = c.makeFifoPipe(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	var cmd *exec.Cmd

	if cmd, err = c.makeCmd(executor, pipePath, args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go c.feedPipe(ar, wg, pipePath, tz)
	go c.exec(wg, pipePath, cmd)

	wg.Wait()
	errors.Print("done waiting")

	return
}

func (c Exec) getZettel(
	u *umwelt.Umwelt,
	hString string,
) (
	tz zettel_transacted.Zettel,
	ar io.ReadCloser,
	executor konfig.RemoteScript,
	err error,
) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(u); err != nil {
		err = errors.Wrap(err)
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
		err = errors.Wrap(err)
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
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Exec) makeFifoPipe(tz zettel_transacted.Zettel) (p string, err error) {
	h := tz.Named.Hinweis
	var d string

	if d, err = os.MkdirTemp("", h.Kopf()); err != nil {
		err = errors.Wrap(err)
		return
	}

	p = path.Join(d, h.Schwanz()+"."+tz.Named.Stored.Zettel.Typ.String())

	if err = syscall.Mknod(p, syscall.S_IFIFO|0666, 0); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.Remove(p); err != nil {
		err = errors.Wrap(err)
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
		err = errors.Wrap(err)
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
	tz zettel_transacted.Zettel,
) (err error) {
	defer wg.Done()
	var pipeFileWriter *os.File

	if pipeFileWriter, err = os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer pipeFileWriter.Close()

	defer ar.Close()

	if _, err = io.Copy(pipeFileWriter, ar); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Print("done copying")

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

	errors.Print("start running")

	if err = cmd.Run(); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Print("done running")

	return
}

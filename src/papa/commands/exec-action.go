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
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/konfig"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/typ"
	"github.com/friedenberg/zit/src/foxtrot/id_set"
	"github.com/friedenberg/zit/src/hotel/zettel_named"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type ExecAction struct {
	Action collections.StringValue
}

func init() {
	registerCommand(
		"exec-action",
		func(f *flag.FlagSet) Command {
			c := &ExecAction{}

			f.Var(&c.Action, "action", "which Typ action to execute")

			return commandWithIds{
				CommandWithIds: c,
			}
		},
	)
}

func (c ExecAction) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet(
		id_set.ProtoId{
			MutableId: &konfig.Id{},
		},
		id_set.ProtoId{
			MutableId: &hinweis.Hinweis{},
			Expand: func(v string) (out string, err error) {
				var h hinweis.Hinweis
				h, err = u.StoreObjekten().ExpandHinweisString(v)
				out = h.String()
				return
			},
		},
		id_set.ProtoId{
			MutableId: &etikett.Etikett{},
			Expand: func(v string) (out string, err error) {
				var e etikett.Etikett
				e, err = u.StoreObjekten().ExpandEtikettString(v)
				out = e.String()
				return
			},
		},
		id_set.ProtoId{
			MutableId: &typ.Typ{},
		},
		id_set.ProtoId{
			MutableId: &ts.Time{},
		},
	)

	return
}

func (c ExecAction) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	if !c.Action.WasSet() {
		err = errors.Normal(errors.Errorf("Action must be provided"))
		return
	}

	query := zettel_transacted.WriterIds(
		zettel_named.FilterIdSet{
			Set: ids,
			// Or:  c.Or,
		},
	)

	iter := func(tz *zettel_transacted.Zettel) (err error) {
		var executor konfig.RemoteScript
		var ar io.ReadCloser

		typ := tz.Named.Stored.Zettel.Typ.String()

		typKonfig := u.Konfig().GetTyp(typ)

		if typKonfig == nil {
			err = errors.Normal(errors.Errorf("Typ does not have an exec-command set: %s", typ))
			return
		}

		ok := false
		executor, ok = typKonfig.Actions[c.Action.String()]

		if !ok {
			err = errors.Normal(errors.Errorf("Typ does not have action: %s", c.Action))
			return
		}

		if ar, err = u.StoreObjekten().AkteReader(tz.Named.Stored.Zettel.Akte); err != nil {
			err = errors.Wrap(err)
			return
		}

		var pipePath string

		if pipePath, err = c.makeFifoPipe(tz); err != nil {
			err = errors.Wrap(err)
			return
		}

		var cmd *exec.Cmd

		if cmd, err = c.makeCmd(executor, pipePath); err != nil {
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

	if err = u.StoreObjekten().ReadAllSchwanzenTransacted(
		query.WriteZettelTransacted,
		iter,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c ExecAction) makeFifoPipe(tz *zettel_transacted.Zettel) (p string, err error) {
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

func (c ExecAction) makeCmd(
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

func (c ExecAction) feedPipe(
	ar io.ReadCloser,
	wg *sync.WaitGroup,
	p string,
	tz *zettel_transacted.Zettel,
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

func (c ExecAction) exec(
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

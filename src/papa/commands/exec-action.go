package commands

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

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
		typ := tz.Named.Stored.Zettel.Typ.String()

		typKonfig := u.Konfig().GetTyp(typ)

		if typKonfig == nil {
			err = errors.Normal(errors.Errorf("Typ does not have an exec-command set: %s", typ))
			return
		}

		executor, ok := typKonfig.Actions[c.Action.String()]

		if !ok {
			err = errors.Normal(errors.Errorf("Typ does not have action: %s", c.Action))
			return
		}

		if err = c.runExecutor(u, executor, tz); err != nil {
			err = errors.Wrap(err)
			return
		}

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

func (c ExecAction) runExecutor(
	u *umwelt.Umwelt,
	executor *konfig.ScriptConfig,
	z *zettel_transacted.Zettel,
) (err error) {
	var cmd *exec.Cmd

	if cmd, err = executor.Cmd(); err != nil {
		err = errors.Wrap(err)
		return
	}

	env := map[string]string{
		"ZETTEL": z.Named.Hinweis.String(),
    "ZIT_BIN": u.Standort().Executable(),
	}

	envCollapsed := make([]string, 0, len(env))

	for k, v := range env {
		envCollapsed = append(envCollapsed, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Env = envCollapsed

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	var wc io.WriteCloser

	if wc, err = cmd.StdinPipe(); err != nil {
		err = errors.Wrap(err)
		return
	}

	chDone := make(chan struct{})

	go func() {
		defer func() {
			chDone <- struct{}{}
		}()

		defer errors.Deferred(&err, wc.Close)

		var ar io.ReadCloser

		if ar, err = u.StoreObjekten().AkteReader(z.Named.Stored.Zettel.Akte); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, ar.Close)

		if _, err = io.Copy(wc, ar); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = wc.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}()

	if err = cmd.Start(); err != nil {
		err = errors.Wrap(err)
		return
	}

	<-chDone

	if err = cmd.Wait(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

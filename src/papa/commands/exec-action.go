package commands

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/id_set"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/typ_toml"
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
			MutableId: &kennung.Konfig{},
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
			MutableId: &kennung.Etikett{},
			Expand: func(v string) (out string, err error) {
				var e kennung.Etikett
				e, err = u.StoreObjekten().ExpandEtikettString(v)
				out = e.String()
				return
			},
		},
		id_set.ProtoId{
			MutableId: &kennung.Typ{},
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
		zettel.FilterIdSet{
			Set: ids,
			// Or:  c.Or,
		},
	)

	iter := func(tz *zettel_transacted.Zettel) (err error) {
		typ := tz.Named.Stored.Objekte.Typ.String()

		typKonfig := u.Konfig().Transacted.Objekte.GetTyp(typ)

		if typKonfig == nil {
			err = errors.Normal(errors.Errorf("Typ does not have an exec-command set: %s", typ))
			return
		}

		executor, ok := typKonfig.Typ.Actions[c.Action.String()]

		if !ok {
			err = errors.Normalf("Typ '%s' does not have action '%s'", typ, c.Action)
			return
		}

		if err = c.runExecutor(u, executor, tz); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = u.StoreObjekten().Zettel().ReadAllSchwanzenTransacted(
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
	executor *typ_toml.Action,
	z *zettel_transacted.Zettel,
) (err error) {
	var cmd *exec.Cmd

	if cmd, err = executor.Cmd(); err != nil {
		err = errors.Wrap(err)
		return
	}

	env := map[string]string{
		"ZIT_ZETTEL": z.Named.Kennung.String(),
		"ZIT_BIN":    u.Standort().Executable(),
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

		if ar, err = u.StoreObjekten().AkteReader(z.Named.Stored.Objekte.Akte); err != nil {
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

package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/hinweis"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/foxtrot/id_set"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type ExecAction struct {
	Action collections.StringValue
}

func init() {
	registerCommand(
		"exec-action",
		func(f *flag.FlagSet) Command {
			c := &ExecAction{}

			f.Var(&c.Action, "action", "which Konfig action to execute")

			return commandWithIds{
				CommandWithIds: c,
			}
		},
	)
}

func (c ExecAction) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet(
		id_set.ProtoId{
			Setter: &kennung.Konfig{},
		},
		id_set.ProtoId{
			Setter: &hinweis.Hinweis{},
			Expand: func(v string) (out string, err error) {
				var h hinweis.Hinweis
				h, err = u.StoreObjekten().Abbr().ExpandHinweisString(v)
				out = h.String()
				return
			},
		},
		id_set.ProtoId{
			Setter: &kennung.Etikett{},
			Expand: func(v string) (out string, err error) {
				var e kennung.Etikett
				e, err = u.StoreObjekten().Abbr().ExpandEtikettString(v)
				out = e.String()
				return
			},
		},
		id_set.ProtoId{
			Setter: &kennung.Typ{},
		},
		id_set.ProtoId{
			Setter: &ts.Time{},
		},
	)

	return
}

func (c ExecAction) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	if !c.Action.WasSet() {
		err = errors.Normal(errors.Errorf("Action must be provided"))
		return
	}

	var sc script_config.ScriptConfig
	ok := false

	if sc, ok = u.Konfig().Actions[c.Action.String()]; !ok {
		err = errors.Normalf(
			"Konfig Action '%s' not found",
			c.Action.String(),
		)

		return
	}

	query := zettel.WriterIds{
		Filter: id_set.Filter{
			Set: ids,
			// Or:  c.Or,
		},
	}

	hinweisen := hinweis.MakeMutableSet()

	if err = u.StoreWorkingDirectory().ReadMany(
		collections.MakeChain(
			query.WriteZettelTransacted,
			func(z *zettel.Transacted) (err error) {
				return hinweisen.Add(z.Sku.Kennung)
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.runAction(
		u,
		sc,
		hinweisen,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c ExecAction) runAction(
	u *umwelt.Umwelt,
	sc script_config.ScriptConfig,
	hinweisen hinweis.MutableSet,
) (err error) {
	var wt io.WriterTo

	if wt, err = script_config.MakeWriterTo(
		sc,
		map[string]string{
			"ZIT_BIN": u.Standort().Executable(),
		},
		hinweisen.Strings()...,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = wt.WriteTo(u.Out()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

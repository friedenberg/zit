package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/juliett/zettel"
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

			f.Var(&c.Action, "action", "which Konfig action to execute")

			return commandWithIds{
				CommandWithIds: c,
			}
		},
	)
}

func (c ExecAction) ProtoIdSet(u *umwelt.Umwelt) (is kennung.ProtoIdSet) {
	is = kennung.MakeProtoIdSet(
		kennung.ProtoId{
			Setter: &kennung.Konfig{},
		},
		kennung.ProtoId{
			Setter: &kennung.Hinweis{},
			Expand: func(v string) (out string, err error) {
				var h kennung.Hinweis
				h, err = u.StoreObjekten().GetAbbrStore().ExpandHinweisString(v)
				out = h.String()
				return
			},
		},
		kennung.ProtoId{
			Setter: &kennung.Etikett{},
			Expand: func(v string) (out string, err error) {
				var e kennung.Etikett
				e, err = u.StoreObjekten().GetAbbrStore().ExpandEtikettString(v)
				out = e.String()
				return
			},
		},
		kennung.ProtoId{
			Setter: &kennung.Typ{},
		},
		kennung.ProtoId{
			Setter: &ts.Time{},
		},
	)

	return
}

func (c ExecAction) RunWithIds(u *umwelt.Umwelt, ids kennung.Set) (err error) {
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
		Filter: kennung.Filter{
			Set: ids,
			// Or:  c.Or,
		},
	}

	hinweisen := kennung.MakeHinweisMutableSet()

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
	hinweisen kennung.HinweisMutableSet,
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

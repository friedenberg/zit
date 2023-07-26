package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type ExecAction struct {
	Action values.String
}

func init() {
	registerCommandWithQuery(
		"exec-action",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &ExecAction{}

			f.Var(&c.Action, "action", "which Konfig action to execute")

			return c
		},
	)
}

func (c ExecAction) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
	)
}

func (c ExecAction) RunWithQuery(
	u *umwelt.Umwelt,
	ms kennung.MetaSet,
) (err error) {
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

	zids, ok := ms.Get(gattung.Zettel)

	if !ok {
		return
	}

	hinweisen := kennung.MakeHinweisMutableSet()

	if err = u.StoreObjekten().Zettel().Query(
		zids,
		func(z *zettel.Transacted) (err error) {
			return hinweisen.Add(z.Sku.GetKennung())
		},
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
		iter.Strings[kennung.Hinweis](hinweisen)...,
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

package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/values"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/src/delta/gattungen"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/matcher"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
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
	ms matcher.Query,
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

	hinweisen := collections_value.MakeMutableValueSet[kennung.Kennung](nil)

	if err = u.StoreObjekten().QueryWithCwd(
		ms,
		func(z *sku.Transacted) (err error) {
			return hinweisen.Add(z.GetKennung())
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
	hinweisen schnittstellen.SetLike[kennung.Kennung],
) (err error) {
	var wt io.WriterTo

	if wt, err = script_config.MakeWriterTo(
		sc,
		map[string]string{
			"ZIT_BIN": u.Standort().Executable(),
		},
		iter.Strings[kennung.Kennung](hinweisen)...,
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

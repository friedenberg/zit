package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
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

func (c ExecAction) DefaultGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
	)
}

func (c ExecAction) RunWithQuery(
	u *env.Env,
	ms *query.Group,
) (err error) {
	if !c.Action.WasSet() {
		err = errors.Normal(errors.Errorf("Action must be provided"))
		return
	}

	var sc script_config.ScriptConfig
	ok := false

	if sc, ok = u.GetConfig().Actions[c.Action.String()]; !ok {
		err = errors.Normalf(
			"Konfig Action '%s' not found",
			c.Action.String(),
		)

		return
	}

	object_id_provider := collections_value.MakeMutableValueSet[ids.IdLike](nil)

	if err = u.GetStore().QueryTransacted(
		ms,
		func(z *sku.Transacted) (err error) {
			return object_id_provider.Add(z.GetObjectId())
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.runAction(
		u,
		sc,
		object_id_provider,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c ExecAction) runAction(
	u *env.Env,
	sc script_config.ScriptConfig,
	object_id_provider interfaces.SetLike[ids.IdLike],
) (err error) {
	var wt io.WriterTo

	if wt, err = script_config.MakeWriterTo(
		sc,
		map[string]string{
			"ZIT_BIN": u.GetFSHome().Executable(),
		},
		iter.Strings[ids.IdLike](object_id_provider)...,
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

package commands

import (
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type CommandWithQuery interface {
	RunWithQuery(store *umwelt.Umwelt, ids *query.Group) error
	DefaultGattungen() kennung.Gattung
}

type commandWithQuery struct {
	CommandWithQuery
}

type CompletionGattungGetter interface {
	CompletionGattung() kennung.Gattung
}

func (c commandWithQuery) Complete(
	u *umwelt.Umwelt,
	args ...string,
) (err error) {
	var cgg CompletionGattungGetter
	ok := false

	if cgg, ok = c.CommandWithQuery.(CompletionGattungGetter); !ok {
		return
	}

	w := sku_fmt.MakeWriterComplete(os.Stdout)
	defer errors.DeferredCloser(&err, w)

	b := u.MakeQueryBuilderExcludingHidden(cgg.CompletionGattung())

	var qg *query.Group

	if qg, err = b.BuildQueryGroup(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetStore().QueryWithoutCwd(
		qg,
		w.WriteOne,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c commandWithQuery) Run(u *umwelt.Umwelt, args ...string) (err error) {
	builder := u.MakeQueryBuilderExcludingHidden(c.DefaultGattungen())

	type withDefaultSigil interface {
		DefaultSigil() kennung.Sigil
	}

	if wds, ok := c.CommandWithQuery.(withDefaultSigil); ok {
		builder.WithDefaultSigil(wds.DefaultSigil())
	}

	var ids *query.Group

	if ids, err = builder.BuildQueryGroup(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithQuery(u, ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

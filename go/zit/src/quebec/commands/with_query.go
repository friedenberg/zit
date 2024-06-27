package commands

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type CommandWithQuery interface {
	RunWithQuery(store *umwelt.Umwelt, ids *query.Group) error
}

type commandWithQuery struct {
	CommandWithQuery
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

	if err = u.GetStore().Query(
		qg,
		w.WriteOne,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c commandWithQuery) Run(u *umwelt.Umwelt, args ...string) (err error) {
	b := u.MakeQueryBuilderExcludingHidden(kennung.MakeGattung())

	if dgg, ok := c.CommandWithQuery.(DefaultGattungGetter); ok {
		b = b.WithDefaultGattungen(dgg.DefaultGattungen())
	}

	if dsg, ok := c.CommandWithQuery.(DefaultSigilGetter); ok {
		b.WithDefaultSigil(dsg.DefaultSigil())
	}

	if qbm, ok := c.CommandWithQuery.(QueryBuilderModifier); ok {
		qbm.ModifyBuilder(b)
	}

	var qg *query.Group

	if qg, err = b.BuildQueryGroup(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithQuery(u, qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

package commands

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type CommandWithExternalQuery interface {
	RunWithExternalQuery(
		store *umwelt.Umwelt,
		ids sku.ExternalQueryWithKasten,
	) error
}

type commandWithExternalQuery struct {
	CommandWithExternalQuery
	sku.ExternalQueryWithKasten
}

func (c commandWithExternalQuery) Complete(
	u *umwelt.Umwelt,
	args ...string,
) (err error) {
	var cgg CompletionGattungGetter
	ok := false

	if cgg, ok = c.CommandWithExternalQuery.(CompletionGattungGetter); !ok {
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

	if err = u.GetStore().QueryOld(
		qg,
		w.WriteOne,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c commandWithExternalQuery) Run(u *umwelt.Umwelt, args ...string) (err error) {
	b := u.MakeQueryBuilderExcludingHidden(kennung.MakeGattung())

	if dgg, ok := c.CommandWithExternalQuery.(DefaultGattungGetter); ok {
		b = b.WithDefaultGattungen(dgg.DefaultGattungen())
	}

	if dsg, ok := c.CommandWithExternalQuery.(DefaultSigilGetter); ok {
		b.WithDefaultSigil(dsg.DefaultSigil())
	}

	if qbm, ok := c.CommandWithExternalQuery.(QueryBuilderModifier); ok {
		qbm.ModifyBuilder(b)
	}

	if c.Queryable, err = b.BuildQueryGroup(
		args...,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithExternalQuery(u, c.ExternalQueryWithKasten); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

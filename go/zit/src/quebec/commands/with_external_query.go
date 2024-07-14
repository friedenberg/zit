package commands

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type CommandWithQuery interface {
	RunWithQuery(
		store *umwelt.Umwelt,
		ids *query.Group,
	) error
}

type commandWithQuery struct {
	CommandWithQuery
	sku.ExternalQueryOptions
	*query.Group
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

	if c.Group, err = b.BuildQueryGroupWithRepoId(
		c.RepoId,
		c.ExternalQueryOptions,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetStore().QueryWithKasten(
		c.Group,
		w.WriteOne,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c commandWithQuery) Run(u *umwelt.Umwelt, args ...string) (err error) {
	b := u.MakeQueryBuilderExcludingHidden(ids.MakeGenre())

	if dgg, ok := c.CommandWithQuery.(DefaultGattungGetter); ok {
		b = b.WithDefaultGenres(dgg.DefaultGattungen())
	}

	if dsg, ok := c.CommandWithQuery.(DefaultSigilGetter); ok {
		b.WithDefaultSigil(dsg.DefaultSigil())
	}

	if qbm, ok := c.CommandWithQuery.(QueryBuilderModifier); ok {
		qbm.ModifyBuilder(b)
	}

	if c.Group, err = b.BuildQueryGroupWithRepoId(
		c.RepoId,
		c.ExternalQueryOptions,
		args...,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.Group.ExternalQueryOptions = c.ExternalQueryOptions

	if err = c.RunWithQuery(u, c.Group); err != nil {
		ui.Debug().Printf("%#v", err)
		err = errors.Wrap(err)
		return
	}

	return
}

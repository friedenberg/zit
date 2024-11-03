package commands

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type CommandWithQuery interface {
	RunWithQuery(
		store *env.Env,
		ids *query.Group,
	) error
}

type CommandWithQuery2 interface {
	RunWithQuery(
		store *env.Env,
		ids *query.Group,
	) Result
}

type commandWithQuery struct {
	CommandWithQuery
	sku.ExternalQueryOptions
	*query.Group
}

func (c commandWithQuery) Complete(
	u *env.Env,
	args ...string,
) (err error) {
	var cgg CompletionGenresGetter
	ok := false

	if cgg, ok = c.CommandWithQuery.(CompletionGenresGetter); !ok {
		return
	}

	w := sku_fmt.MakeWriterComplete(os.Stdout)
	defer errors.DeferredCloser(&err, w)

	b := u.MakeQueryBuilderExcludingHidden(cgg.CompletionGenres())

	if c.Group, err = b.BuildQueryGroupWithRepoId(
		c.RepoId,
		c.ExternalQueryOptions,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetStore().QueryTransacted(
		c.Group,
		w.WriteOne,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c commandWithQuery) Run(u *env.Env, args ...string) (err error) {
	b := u.MakeQueryBuilderExcludingHidden(ids.MakeGenre())

	if dgg, ok := c.CommandWithQuery.(DefaultGenresGetter); ok {
		b = b.WithDefaultGenres(dgg.DefaultGenres())
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
		err = errors.Wrap(err)
		return
	}

	return
}

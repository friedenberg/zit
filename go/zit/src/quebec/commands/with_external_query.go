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
		store *env.Local,
		ids *query.Group,
	) error
}

type CommandWithQuery2 interface {
	RunWithQuery(
		store *env.Local,
		ids *query.Group,
	) Result
}

type commandWithQuery struct {
	CommandWithQuery
	sku.ExternalQueryOptions
	*query.Group
}

type CompletionGenresGetter interface {
	CompletionGenres() ids.Genre
}

func (c commandWithQuery) Complete(
	u *env.Local,
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
		w.WriteOneTransacted,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c commandWithQuery) Run(u *env.Local, args ...string) (err error) {
	if c.Group, err = u.MakeQueryGroup(
		c.CommandWithQuery,
		c.RepoId,
		c.ExternalQueryOptions,
		args...,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithQuery(u, c.Group); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

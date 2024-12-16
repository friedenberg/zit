package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

type QueryBuilderModifier interface {
	ModifyBuilder(*query.Builder)
}

type DefaultSigilGetter interface {
	DefaultSigil() ids.Sigil
}

type DefaultGenresGetter interface {
	DefaultGenres() ids.Genre
}

type CompletionGenresGetter interface {
	CompletionGenres() ids.Genre
}

func (u *Local) MakeQueryGroup(
	metaBuilder any,
	repoId ids.RepoId,
	externalQueryOptions sku.ExternalQueryOptions,
	args ...string,
) (qg *query.Group, err error) {
	b := u.MakeQueryBuilderExcludingHidden(ids.MakeGenre())

	if dgg, ok := metaBuilder.(DefaultGenresGetter); ok {
		b = b.WithDefaultGenres(dgg.DefaultGenres())
	}

	if dsg, ok := metaBuilder.(DefaultSigilGetter); ok {
		b.WithDefaultSigil(dsg.DefaultSigil())
	}

	if qbm, ok := metaBuilder.(QueryBuilderModifier); ok {
		qbm.ModifyBuilder(b)
	}

	if qg, err = b.BuildQueryGroupWithRepoId(
		repoId,
		externalQueryOptions,
		args...,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	qg.ExternalQueryOptions = externalQueryOptions

	return
}

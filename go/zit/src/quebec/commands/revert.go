package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Revert struct {
	Last bool
}

func init() {
	registerCommandWithQuery(
		"revert",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Revert{}

			f.BoolVar(&c.Last, "last", false, "revert the last changes")

			return c
		},
	)
}

func (c Revert) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
		genres.Tag,
		genres.Type,
		genres.InventoryList,
		genres.Repo,
	)
}

func (c Revert) DefaultGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
		genres.Tag,
		genres.Type,
		// gattung.Bestandsaufnahme,
		genres.Repo,
	)
}

type revertTuple struct {
	*sku.Transacted
	*sha.Sha
}

func (c Revert) RunWithQuery(
	u *env.Env,
	ms *query.Group,
) (err error) {
	f := func(rt revertTuple) (err error) {
		if rt.IsNull() {
			return
		}

		if err = u.GetStore().RevertTo(rt.Transacted, rt.Sha); err != nil {
			err = errors.Wrapf(err, "Sha %s", rt.Sha)
			return
		}

		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	switch {
	case c.Last:
		err = c.runRevertFromLast(u, f)

	default:
		err = c.runRevertFromQuery(u, ms, f)
	}

	if err != nil {
		return
	}

	return
}

func (c Revert) runRevertFromQuery(
	u *env.Env,
	eq *query.Group,
	f interfaces.FuncIter[revertTuple],
) (err error) {
	if err = u.GetStore().QueryWithKasten(
		eq,
		func(z *sku.Transacted) (err error) {
			return f(revertTuple{Transacted: z, Sha: z.Metadata.Mutter()})
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Revert) runRevertFromLast(
	u *env.Env,
	f interfaces.FuncIter[revertTuple],
) (err error) {
	s := u.GetStore()

	var b *sku.Transacted

	if b, err = s.GetBestandsaufnahmeStore().ReadLast(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetBestandsaufnahmeStore().StreamInventoryList(
		b.GetBlobSha(),
		func(sk *sku.Transacted) (err error) {
			return f(revertTuple{Transacted: sk, Sha: sk.Metadata.Mutter()})
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
)

type Fsck struct {
	Genres ids.Genre
}

func init() {
	registerCommandWithQuery(
		"fsck",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Fsck{
				Genres: ids.MakeGenre(genres.Tag, genres.Type, genres.Zettel),
			}

			f.Var(&c.Genres, "genres", "")

			return c
		},
	)
}

func (c Fsck) RunWithQuery(
	u *read_write_repo_local.Repo,
	qg *query.Group,
) {
	p := u.PrinterTransacted()

	if err := u.GetStore().QueryTransacted(
		qg,
		func(sk *sku.Transacted) (err error) {
			if !c.Genres.Contains(sk.GetGenre()) {
				return
			}

			blobSha := sk.GetBlobSha()

			if u.GetRepoLayout().HasBlob(blobSha) {
				return
			}

			if err = p(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		u.CancelWithError(err)
	}
}

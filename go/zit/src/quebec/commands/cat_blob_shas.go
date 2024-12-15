package commands

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type CatBlobShas struct{}

func init() {
	registerCommand(
		"cat-blob-shas",
		func(f *flag.FlagSet) Command {
			c := &CatBlobShas{}

			return c
		},
	)
}

func (c CatBlobShas) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Blob,
	)
}

func (c CatBlobShas) Run(u *env.Local, _ ...string) (err error) {
	if err = u.GetDirectoryLayout().ReadAllShasForGenre(
		genres.Blob,
		func(s *sha.Sha) (err error) {
			_, err = fmt.Fprintln(u.Out(), s)
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

package commands

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type CatBlobShas struct{}

func init() {
	registerCommand(
		"cat-blob-shas",
		func(f *flag.FlagSet) CommandWithLocalWorkingCopy {
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

func (c CatBlobShas) RunWithLocalWorkingCopy(u *local_working_copy.Repo, _ ...string) {
	if err := u.GetRepoLayout().ReadAllShasForGenre(
		genres.Blob,
		func(s *sha.Sha) (err error) {
			_, err = fmt.Fprintln(u.GetUIFile(), s)
			return
		},
	); err != nil {
		u.CancelWithError(err)
	}
}

package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Fsck struct {
	Genres ids.Genre
}

func init() {
	registerCommand(
		"fsck",
		func(f *flag.FlagSet) CommandWithResult {
			c := &Fsck{
				Genres: ids.MakeGenre(genres.Tag, genres.Type, genres.Zettel),
			}

			f.Var(&c.Genres, "genres", "")

			return c
		},
	)
}

func (c Fsck) Run(u *env.Local, args ...string) (result Result) {
	p := u.PrinterTransacted()

	if result.Error = u.GetStore().QueryPrimitive(
		sku.MakePrimitiveQueryGroup(),
		func(sk *sku.Transacted) (err error) {
			if !c.Genres.Contains(sk.GetGenre()) {
				return
			}

			blobSha := sk.GetBlobSha()

			if u.GetDirectoryLayout().HasBlob(blobSha) {
				return
			}

			if err = p(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); result.Error != nil {
		result.Error = errors.Wrap(result.Error)
		return
	}

	return
}

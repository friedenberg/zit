package commands

import (
	"bufio"
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Export struct {
	AgeIdentity     age.Identity
	CompressionType immutable_config.CompressionType
}

func init() {
	registerCommandWithQuery(
		"export",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Export{
				CompressionType: immutable_config.CompressionTypeEmpty,
			}

			f.Var(&c.AgeIdentity, "age-identity", "")
			c.CompressionType.AddToFlagSet(f)

			return c
		},
	)
}

func (c Export) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c Export) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
}

func (c Export) RunWithQuery(u *repo_local.Repo, qg *query.Group) {
	var list *sku.List

	{
		var err error

		if list, err = u.MakeInventoryList(qg); err != nil {
			u.CancelWithError(err)
		}
	}

	var ag age.Age

	if err := ag.AddIdentity(c.AgeIdentity); err != nil {
		u.CancelWithErrorAndFormat(err, "age-identity: %q", &c.AgeIdentity)
	}

	var wc io.WriteCloser

	o := repo_layout.WriteOptions{
		Age:             &ag,
		CompressionType: c.CompressionType,
		Writer:          u.GetUIFile(),
	}

	{
		var err error

		if wc, err = repo_layout.NewWriter(o); err != nil {
			u.CancelWithError(err)
		}
	}

	defer u.MustClose(wc)

	bw := bufio.NewWriter(wc)
	defer u.MustFlush(bw)

	printer := u.MakePrinterBoxArchive(bw, u.GetConfig().PrintOptions.PrintTime)

	var sk *sku.Transacted
	var hasMore bool

	for {
		u.ContinueOrPanicOnDone()

		sk, hasMore = list.Pop()

		if !hasMore {
			break
		}

		if err := printer(sk); err != nil {
			u.CancelWithError(err)
		}
	}
}

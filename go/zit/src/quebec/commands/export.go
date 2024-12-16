package commands

import (
	"bufio"
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
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

func (c Export) RunWithQuery(u *env.Local, qg *query.Group) (err error) {
	var list *sku.List

	if list, err = u.MakeInventoryList(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ag age.Age

	if err = ag.AddIdentity(c.AgeIdentity); err != nil {
		err = errors.Wrapf(err, "age-identity: %q", &c.AgeIdentity)
		return
	}

	var wc io.WriteCloser

	o := dir_layout.WriteOptions{
		Age:             &ag,
		CompressionType: c.CompressionType,
		Writer:          u.Out(),
	}

	if wc, err = dir_layout.NewWriter(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, wc)

	bw := bufio.NewWriter(wc)
	defer errors.DeferredFlusher(&err, bw)

	printer := u.MakePrinterBoxArchive(bw, u.GetConfig().PrintOptions.PrintTime)

	var sk *sku.Transacted
	var hasMore bool

	for {
		sk, hasMore = list.Pop()

		if !hasMore {
			break
		}

		if err = printer(sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

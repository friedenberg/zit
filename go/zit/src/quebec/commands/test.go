package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type Test struct{}

func init() {
	registerCommandWithQuery(
		"test",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Test{}

			return c
		},
	)
}

func (c Test) CompletionGattung() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
		genres.Tag,
		genres.Type,
		genres.InventoryList,
		genres.Repo,
	)
}

func (c Test) DefaultGattungen() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
		genres.Tag,
		genres.Type,
		// gattung.Bestandsaufnahme,
		genres.Repo,
	)
}

func (c Test) RunWithQuery(
	u *umwelt.Umwelt,
	ms *query.Group,
) (err error) {
	if err = u.GetStore().QueryWithKasten(
		ms,
		iter.MakeSyncSerializer(
			func(o *sku.Transacted) (err error) {
				var sk *sku.Transacted

				if sk, err = u.GetStore().GetVerzeichnisse().ReadOneObjectId(
					&o.Kennung,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				defer sku.GetTransactedPool().Put(sk)

				if object_metadata.EqualerSansTai.Equals(o.GetMetadata(), sk.GetMetadata()) {
					return
				}

				ui.Out().Print(o.GetObjectId())
				ui.Debug().Print(o.GetTai())
				ui.Debug().Print(sk.GetTai())

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/query"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
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

func (c Test) CompletionGattung() kennung.Gattung {
	return kennung.MakeGattung(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
		gattung.Bestandsaufnahme,
		gattung.Kasten,
	)
}

func (c Test) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
		// gattung.Bestandsaufnahme,
		gattung.Kasten,
	)
}

func (c Test) RunWithQuery(u *umwelt.Umwelt, ms *query.Group) (err error) {
	if err = u.GetStore().QueryWithCwd(
		ms,
		iter.MakeSyncSerializer(
			func(o *sku.Transacted) (err error) {
				var sk *sku.Transacted

				if sk, err = u.GetStore().GetVerzeichnisse().ReadOneKennung(
					&o.Kennung,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				defer sku.GetTransactedPool().Put(sk)

				if metadatei.EqualerSansTai.Equals(o.GetMetadatei(), sk.GetMetadatei()) {
					return
				}

				log.Out().Print(o.GetKennung())
				log.Debug().Print(o.GetTai())
				log.Debug().Print(sk.GetTai())

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

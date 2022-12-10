package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/oscar/store_fs"
	"github.com/friedenberg/zit/src/papa/umwelt"
)

type Checkout struct {
	store_fs.CheckoutMode
	Or    bool
	Force bool
}

func init() {
	registerCommand(
		"checkout",
		func(f *flag.FlagSet) Command {
			c := &Checkout{}

			f.BoolVar(&c.Or, "or", false, "allow optional criteria instead of required")
			f.Var(&c.CheckoutMode, "mode", "mode for checking out the zettel")
			f.BoolVar(&c.Force, "force", false, "force update checked out zettels, even if they will overwrite existing checkouts")

			return commandWithIds{
				CommandWithIds: c,
			}
		},
	)
}

func (c Checkout) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet(
		id_set.ProtoId{
			MutableId: &hinweis.Hinweis{},
			Expand: func(v string) (out string, err error) {
				var h hinweis.Hinweis
				h, err = u.StoreObjekten().Abbr().ExpandHinweisString(v)
				out = h.String()
				return
			},
		},
		// id_set.ProtoId{
		// 	MutableId: &sha.Sha{},
		// },
		id_set.ProtoId{
			MutableId: &kennung.Etikett{},
		},
		id_set.ProtoId{
			MutableId: &kennung.Typ{},
		},
		id_set.ProtoId{
			MutableId: &ts.Time{},
		},
	)

	return
}

func (c Checkout) RunWithIds(s *umwelt.Umwelt, ids id_set.Set) (err error) {
	options := store_fs.CheckoutOptions{
		CheckoutMode: c.CheckoutMode,
		Format:       zettel.Text{},
	}

	query := zettel.WriterIds{
		Filter: id_set.Filter{
			Set: ids,
			Or:  c.Or,
		},
	}

	if _, err = s.StoreWorkingDirectory().Checkout(
		options,
		query.WriteZettelTransacted,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/charlie/typ"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/kilo/umwelt"
	"github.com/friedenberg/zit/src/lima/user_ops"
)

type Checkout struct {
	store_working_directory.CheckoutMode
	Force bool
}

func init() {
	registerCommand(
		"checkout",
		func(f *flag.FlagSet) Command {
			c := &Checkout{}

			f.Var(&c.CheckoutMode, "mode", "mode for checking out the zettel")
			f.BoolVar(&c.Force, "force", false, "force update checked out zettels, even if they will overwrite existing checkouts")

			return commandWithIds{
				CommandWithIds: c,
			}
		},
	)
}

func (c Checkout) ProtoIdList(u *umwelt.Umwelt) (is id_set.ProtoIdList) {
	is = id_set.MakeProtoIdList(
		id_set.ProtoId{
			MutableId: &hinweis.Hinweis{},
			Expand: func(v string) (out string, err error) {
				var h hinweis.Hinweis
				h, err = u.StoreObjekten().ExpandHinweisString(v)
				out = h.String()
				return
			},
		},
		// id_set.ProtoId{
		// 	MutableId: &sha.Sha{},
		// },
		id_set.ProtoId{
			MutableId: &etikett.Etikett{},
		},
		id_set.ProtoId{
			MutableId: &typ.Typ{},
		},
		id_set.ProtoId{
			MutableId: &ts.Time{},
		},
	)

	return
}

func (c Checkout) RunWithIds(s *umwelt.Umwelt, ids id_set.Set) (err error) {
	checkoutOp := user_ops.Checkout{
		Umwelt: s,
		CheckoutOptions: store_working_directory.CheckoutOptions{
			CheckoutMode: c.CheckoutMode,
			Format:       zettel.Text{},
		},
	}

	if _, err = checkoutOp.RunMany(ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

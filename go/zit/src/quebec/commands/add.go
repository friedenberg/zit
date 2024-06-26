package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Add struct {
	Dedupe              bool
	Delete              bool
	OpenAkten           bool
	CheckoutAktenAndRun string
	Organize            bool
	Filter              script_value.ScriptValue

	zettel.ProtoZettel
}

func init() {
	registerCommandWithQuery(
		"add",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Add{
				ProtoZettel: zettel.MakeEmptyProtoZettel(),
			}

			f.BoolVar(
				&c.Dedupe,
				"dedupe",
				false,
				"deduplicate added Zettelen based on Akte sha",
			)

			f.BoolVar(
				&c.Delete,
				"delete",
				false,
				"delete the zettel and akte after successful checkin",
			)

			f.BoolVar(&c.OpenAkten, "open-akten", false, "also open the Akten")

			f.StringVar(
				&c.CheckoutAktenAndRun,
				"each-akte",
				"",
				"checkout each Akte and run a utility",
			)

			f.BoolVar(&c.Organize, "organize", false, "")

			c.AddToFlagSet(f)

			errors.TodoP2(
				"add support for restricted query to specific gattung",
			)
			return c
		},
	)
}

func (c Add) ModifyBuilder(b *query.Builder) {
	b.WithDefaultGattungen(kennung.MakeGattung(gattung.Zettel)).
		WithDoNotMatchEmpty()
}

func (c Add) RunWithQuery(
	u *umwelt.Umwelt,
	qg *query.Group,
) (err error) {
	zettelsFromAkteOp := user_ops.ZettelFromExternalAkte{
		Umwelt:      u,
		ProtoZettel: c.ProtoZettel,
		Filter:      c.Filter,
		Delete:      c.Delete,
		Dedupe:      c.Dedupe,
	}

	var zettelsFromAkteResults sku.TransactedMutableSet

	if zettelsFromAkteResults, err = zettelsFromAkteOp.Run(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.openAktenIfNecessary(u, zettelsFromAkteResults); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !c.Organize {
		return
	}

	opOrganize := user_ops.Organize{
		Umwelt:    u,
		Metadatei: c.Metadatei,
	}

	if err = u.GetKonfig().DefaultEtiketten.EachPtr(
		opOrganize.Metadatei.AddEtikettPtr,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = opOrganize.Run(qg, zettelsFromAkteResults); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO [radi/kof !task "add support for kasten in checkouts and external" project-2021-zit-features today zz-inbox]
func (c Add) openAktenIfNecessary(
	u *umwelt.Umwelt,
	zettels sku.TransactedMutableSet,
) (err error) {
	if !c.OpenAkten && c.CheckoutAktenAndRun == "" {
		return
	}

	opCheckout := user_ops.Checkout{
		Umwelt: u,
		Options: checkout_options.Options{
			CheckoutMode: checkout_mode.ModeAkteOnly,
		},
		Utility: c.CheckoutAktenAndRun,
	}

	if _, err = opCheckout.Run(
		zettels,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

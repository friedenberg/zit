package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Diff struct {
	Kasten kennung.Kasten
}

func init() {
	registerCommandWithQuery(
		"diff",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Diff{}

			f.Var(&c.Kasten, "kasten", "none or Chrome")

			return c
		},
	)
}

func (c Diff) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(gattung.TrueGattung()...)
}

func (c Diff) ModifyBuilder(
	b *query.Builder,
) {
	b.WithHidden(nil)
}

func (c Diff) RunWithQuery(
	u *umwelt.Umwelt,
	qg *query.Group,
) (err error) {
	co := checkout_options.TextFormatterOptions{
		DoNotWriteEmptyBezeichnung: true,
	}

	opDiffFS := user_ops.DiffFS{
		Umwelt: u,
		Inline: metadatei.MakeTextFormatterMetadateiInlineAkte(
			co,
			u.Standort(),
			nil,
		),
		Metadatei: metadatei.MakeTextFormatterMetadateiOnly(
			co,
			u.Standort(),
			nil,
		),
	}

	if err = u.GetStore().ReadExternal(
		query.GroupWithKasten{
			Group:  qg,
			Kasten: c.Kasten,
		},
		func(co sku.CheckedOutLike) (err error) {
			switch cot := co.(type) {
			case *store_fs.CheckedOut:
				if err = opDiffFS.Run(cot); err != nil {
					err = errors.Wrap(err)
					return
				}

				// TODO [radi/kof !task "add support for kasten in checkouts and external" project-2021-zit-features today zz-inbox]
			default:
				ui.Err().Printf("unsupportted type: %T, %s", cot, cot)
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

package commands

import (
	"flag"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/cwd"
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
	pz := u.GetStore().GetCwdFiles()

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

	if err = c.openAktenIfNecessary(u, zettelsFromAkteResults, pz); err != nil {
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

	if err = opOrganize.Run(qg, zettelsFromAkteResults); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Add) openAktenIfNecessary(
	u *umwelt.Umwelt,
	zettels sku.TransactedMutableSet,
	cwd *cwd.CwdFiles,
) (err error) {
	if !c.OpenAkten && c.CheckoutAktenAndRun == "" {
		return
	}

	options := checkout_options.Options{
		CheckoutMode: checkout_mode.ModeAkteOnly,
	}

	var filesAkten []string
	var l sync.Mutex

	if err = zettels.Each(
		func(z *sku.Transacted) (err error) {
			var co *sku.CheckedOut

			if co, err = u.GetStore().CheckoutOne(options, z); err != nil {
				err = errors.Wrap(err)
				return
			}

			e := co.External.GetAkteFD().GetPath()

			if e == "" {
				return iter.MakeErrStopIteration()
			}

			l.Lock()
			filesAkten = append(filesAkten, e)
			l.Unlock()

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.OpenAkten {
		openOp := user_ops.OpenFiles{}

		if err = openOp.Run(u, filesAkten...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if c.CheckoutAktenAndRun != "" {
		eachAkteOp := user_ops.EachAkte{}

		if err = eachAkteOp.Run(
			u,
			c.CheckoutAktenAndRun,
			filesAkten...,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

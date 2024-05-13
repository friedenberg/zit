package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/src/november/umwelt"
	"code.linenisgreat.com/zit/src/papa/user_ops"
)

type New struct {
	Edit      bool
	Delete    bool
	Dedupe    bool
	Count     int
	PrintOnly bool
	Filter    script_value.ScriptValue

	zettel.ProtoZettel
}

func init() {
	registerCommand(
		"new",
		func(f *flag.FlagSet) Command {
			c := &New{
				ProtoZettel: zettel.MakeEmptyProtoZettel(),
			}

			f.BoolVar(
				&c.Delete,
				"delete",
				false,
				"delete the zettel and akte after successful checkin",
			)
			f.BoolVar(
				&c.Dedupe,
				"dedupe",
				false,
				"deduplicate added Zettelen based on Akte sha",
			)
			f.BoolVar(
				&c.Edit,
				"edit",
				true,
				"create a new empty zettel and open EDITOR or VISUAL for editing and then commit the resulting changes",
			)
			f.IntVar(
				&c.Count,
				"count",
				1,
				"when creating new empty zettels, how many to create. otherwise ignored",
			)

			f.Var(
				&c.Filter,
				"filter",
				"a script to run for each file to transform it the standard zettel format",
			)

			c.AddToFlagSet(f)

			return c
		},
	)
}

func (c New) ValidateFlagsAndArgs(
	u *umwelt.Umwelt,
	args ...string,
) (err error) {
	if u.Konfig().DryRun && len(args) == 0 {
		err = errors.Errorf(
			"when -dry-run is set, paths to existing zettels must be provided",
		)
		return
	}

	return
}

func (c New) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if err = c.ValidateFlagsAndArgs(u, args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	cotfo := checkout_options.TextFormatterOptions{}

	f := metadatei.TextFormat{
		TextFormatter: metadatei.MakeTextFormatterMetadateiInlineAkte(
			cotfo,
			u.Standort(),
			nil,
		),
		TextParser: metadatei.MakeTextParser(
			u.Standort(),
			nil,
		),
	}

	var zsc sku.CheckedOutMutableSet

	if len(args) == 0 {
		if zsc, err = c.writeNewZettels(u); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		zsc = collections_value.MakeMutableValueSet[*sku.CheckedOut](nil)

		var zts sku.TransactedMutableSet

		if zts, err = c.readExistingFilesAsZettels(u, f, args...); err != nil {
			err = errors.Wrap(err)
			return
		}

		if c.Edit {
			opCheckout := user_ops.Checkout{
				Umwelt: u,
				Options: checkout_options.Options{
					CheckoutMode:         checkout_mode.ModeObjekteAndAkte,
					TextFormatterOptions: cotfo,
				},
			}

			if zsc, err = opCheckout.Run(zts); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if c.Edit {
		opEdit := user_ops.Edit{
			Umwelt: u,
		}

		if err = opEdit.Run(zsc); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c New) readExistingFilesAsZettels(
	u *umwelt.Umwelt,
	f metadatei.TextParser,
	args ...string,
) (zts sku.TransactedMutableSet, err error) {
	opCreateFromPath := user_ops.CreateFromPaths{
		Umwelt:      u,
		TextParser:  f,
		Filter:      c.Filter,
		Delete:      c.Delete,
		Dedupe:      c.Dedupe,
		ProtoZettel: c.ProtoZettel,
	}

	if zts, err = opCreateFromPath.Run(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c New) writeNewZettels(
	u *umwelt.Umwelt,
) (zsc sku.CheckedOutMutableSet, err error) {
	emptyOp := user_ops.WriteNewZettels{
		Umwelt:   u,
		CheckOut: c.Edit,
	}

	u.Konfig().DefaultEtiketten.EachPtr(c.Metadatei.AddEtikettPtr)

	if zsc, err = emptyOp.RunMany(c.ProtoZettel, c.Count); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

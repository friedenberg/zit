package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/charlie/checkout_options"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/india/objekte_collections"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/user_ops"
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

	f := metadatei.TextFormat{
		TextFormatter: metadatei.MakeTextFormatterMetadateiInlineAkte(
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
		var zts sku.TransactedMutableSet

		if zts, err = c.readExistingFilesAsZettels(u, f, args...); err != nil {
			err = errors.Wrap(err)
			return
		}

		if c.Edit {
			options := checkout_options.Options{
				CheckoutMode: checkout_mode.ModeObjekteAndAkte,
			}

			if zsc, err = u.StoreObjekten().Checkout(
				options,
				u.StoreObjekten().MakeReadAllSchwanzen(gattung.Zettel),
				func(sk *sku.Transacted) (err error) {
					if zts.ContainsKey(sk.GetKennungLike().String()) {
						err = collections.MakeErrStopIteration()
						return
					}

					return
				},
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if c.Edit {
		if err = c.editZettels(
			u,
			zsc,
		); err != nil {
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

func (c New) editZettels(
	u *umwelt.Umwelt,
	zsc sku.CheckedOutMutableSet,
) (err error) {
	if !c.Edit {
		errors.Log().Print("edit set to false, not editing")
		return
	}

	var filesZettelen []string

	if filesZettelen, err = objekte_collections.ToSliceFilesZettelen(zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithCursorLocation(2, 3).
			WithFileType("zit-zettel").
			WithInsertMode().
			Build(),
	}

	if _, err = openVimOp.Run(u, filesZettelen...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	checkinOp := user_ops.Checkin{}

	var ms matcher.Query

	if ms, err = matcher.MakeQueryFromCheckedOutSet(zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = checkinOp.Run(u, ms); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

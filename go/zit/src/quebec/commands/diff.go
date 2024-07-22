package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Diff struct{}

// TODO switch to registerCommandWithExternalQuery
func init() {
	registerCommandWithQuery(
		"diff",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Diff{}

			return c
		},
	)
}

func (c Diff) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Diff) ModifyBuilder(
	b *query.Builder,
) {
	b.WithHidden(nil)
}

func (c Diff) RunWithQuery(
	u *env.Env,
	qg *query.Group,
) (err error) {
	co := checkout_options.TextFormatterOptions{
		DoNotWriteEmptyDescription: true,
	}

	opDiffFS := user_ops.Diff{
		Env: u,
		Inline: object_metadata.MakeTextFormatterMetadataInlineBlob(
			co,
			u.GetFSHome(),
			nil,
		),
		Metadata: object_metadata.MakeTextFormatterMetadataOnly(
			co,
			u.GetFSHome(),
			nil,
		),
	}

	if err = u.GetStore().QueryCheckedOut(
		qg,
		func(co sku.CheckedOutLike) (err error) {
			switch cot := co.(type) {
			case *store_fs.CheckedOut:
				if err = opDiffFS.Run(cot); err != nil {
					err = errors.Wrap(err)
					return
				}

			default:
				ui.Err().Printf("unsupported type: %T, %s", cot, cot)
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	registerCommandWithQuery("status", &Status{})
}

type Status struct {
	command_components.LocalWorkingCopyWithQueryGroup
}

func (c Status) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Status) ModifyBuilder(
	b *query.Builder,
) {
	b.WithHidden(nil).
		WithDefaultSigil(ids.SigilExternal)
}

func (cmd Status) Run(dep command.Dep) {
	u, qg := cmd.MakeLocalWorkingCopyAndQueryGroup(
		dep,
		query.MakeBuilderOptions(cmd),
	)

	pcol := u.PrinterCheckedOut(box_format.CheckedOutHeaderState{})

	if err := u.GetStore().QuerySkuType(
		qg,
		func(co sku.SkuType) (err error) {
			if err = pcol(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		u.CancelWithError(err)
	}
}

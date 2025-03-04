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
	command.Register("status", &Status{})
}

type Status struct {
	command_components.LocalWorkingCopyWithQueryGroup
}

func (c Status) ModifyBuilder(
	b *query.Builder,
) {
	b.WithHidden(nil)
}

func (cmd Status) Run(req command.Request) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		req,
		query.BuilderOptionsOld(
			cmd,
			query.BuilderOptionDefaultGenres(genres.All()...),
			query.BuilderOptionDefaultSigil(ids.SigilExternal),
		),
	)

	envWorkspace := localWorkingCopy.GetEnvWorkspace()
	envWorkspace.AssertInWorkspaceOrOfferToCreate(req)

	pcol := localWorkingCopy.PrinterCheckedOut(box_format.CheckedOutHeaderState{})

	if err := localWorkingCopy.GetStore().QuerySkuType(
		queryGroup,
		func(co sku.SkuType) (err error) {
			if err = pcol(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		localWorkingCopy.CancelWithError(err)
	}
}

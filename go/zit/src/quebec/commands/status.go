package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	pkg_query "code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("status", &Status{})
}

type Status struct {
	command_components.LocalWorkingCopyWithQueryGroup
}

func (c Status) ModifyBuilder(
	b *pkg_query.Builder,
) {
	b.WithHidden(nil)
}

func (cmd Status) Run(req command.Request) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(req)
	localWorkingCopy.GetEnvWorkspace().AssertNotTemporary(req)

	query := cmd.MakeQueryIncludingWorkspace(
		req,
		pkg_query.BuilderOptions(
			pkg_query.BuilderOptionDefaultGenres(genres.All()...),
			pkg_query.BuilderOptionDefaultSigil(ids.SigilExternal),
			pkg_query.BuilderOptionsOld(cmd),
		),
		localWorkingCopy,
		req.PopArgs(),
	)

	printer := localWorkingCopy.PrinterCheckedOut(box_format.CheckedOutHeaderState{})

	if err := localWorkingCopy.GetStore().QuerySkuType(
		query,
		func(co sku.SkuType) (err error) {
			if err = printer(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		localWorkingCopy.CancelWithError(err)
	}
}

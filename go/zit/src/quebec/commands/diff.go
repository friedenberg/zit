package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

// TODO switch to registerCommandWithExternalQuery
func init() {
	command.Register("diff", &Diff{})
}

type Diff struct {
	command_components.LocalWorkingCopyWithQueryGroup
}

func (cmd *Diff) SetFlagSet(f *flag.FlagSet) {
	cmd.LocalWorkingCopyWithQueryGroup.SetFlagSet(f)
}

func (c Diff) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.All()...)
}

func (c Diff) ModifyBuilder(
	b *query.Builder,
) {
	b.WithHidden(nil)
}

func (cmd Diff) Run(dep command.Request) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		dep,
		query.MakeBuilderOptions(cmd),
	)

	o := checkout_options.TextFormatterOptions{
		DoNotWriteEmptyDescription: true,
	}

	opDiffFS := user_ops.Diff{
		Repo: localWorkingCopy,
		TextFormatterFamily: object_metadata.MakeTextFormatterFamily(
			object_metadata.Dependencies{
				EnvDir:    localWorkingCopy.GetEnvRepo(),
				BlobStore: localWorkingCopy.GetEnvRepo(),
			},
		),
	}

	if err := localWorkingCopy.GetStore().QuerySkuType(
		queryGroup,
		func(co sku.SkuType) (err error) {
			if err = opDiffFS.Run(co, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		localWorkingCopy.CancelWithError(err)
	}
}

package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register(
		"fsck",
		&Fsck{
			Genres: ids.MakeGenre(genres.Tag, genres.Type, genres.Zettel),
		},
	)
}

type Fsck struct {
	command_components.LocalWorkingCopyWithQueryGroup

	Genres ids.Genre
}

func (cmd *Fsck) SetFlagSet(f *flag.FlagSet) {
	cmd.LocalWorkingCopyWithQueryGroup.SetFlagSet(f)
	f.Var(&cmd.Genres, "genres", "")
}

func (cmd Fsck) Run(dep command.Request) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		dep,
		query.MakeBuilderOptions(cmd),
	)

	p := localWorkingCopy.PrinterTransacted()

	if err := localWorkingCopy.GetStore().QueryTransacted(
		queryGroup,
		func(sk *sku.Transacted) (err error) {
			if !cmd.Genres.Contains(sk.GetGenre()) {
				return
			}

			blobSha := sk.GetBlobSha()

			if localWorkingCopy.GetRepoLayout().HasBlob(blobSha) {
				return
			}

			if err = p(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		localWorkingCopy.CancelWithError(err)
	}
}

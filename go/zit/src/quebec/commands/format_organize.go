package commands

import (
	"flag"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type FormatOrganize struct {
	Flags organize_text.Flags
}

func init() {
	registerCommandOld(
		"format-organize",
		func(f *flag.FlagSet) WithLocalWorkingCopy {
			c := &FormatOrganize{
				Flags: organize_text.MakeFlags(),
			}

			c.Flags.SetFlagSet(f)

			return c
		},
	)
}

func (c *FormatOrganize) Run(u *local_working_copy.Repo, args ...string) {
	c.Flags.Config = u.GetConfig()

	if len(args) != 1 {
		u.CancelWithErrorf("expected exactly one input argument")
	}

	var fdee fd.FD

	if err := fdee.Set(args[0]); err != nil {
		u.CancelWithError(err)
	}

	var r io.Reader

	if fdee.IsStdin() {
		r = os.Stdin
	} else {
		var f *os.File

		{
			var err error

			if f, err = files.Open(args[0]); err != nil {
				u.CancelWithError(err)
			}
		}

		r = f

		defer u.MustClose(f)
	}

	var ot *organize_text.Text

	readOrganizeTextOp := user_ops.ReadOrganizeFile{}

	var repoId ids.RepoId

	{
		var err error

		if ot, err = readOrganizeTextOp.Run(
			u,
			r,
			organize_text.NewMetadata(repoId),
		); err != nil {
			u.CancelWithError(err)
		}
	}

	ot.Options = c.Flags.GetOptionsWithMetadata(
		u.GetConfig().PrintOptions,
		u.SkuFormatBoxCheckedOutNoColor(),
		u.GetStore().GetAbbrStore().GetAbbr(),
		u.GetExternalLikePoolForRepoId(repoId),
		ot.Metadata,
	)

	if err := ot.Refine(); err != nil {
		u.CancelWithError(err)
	}

	if _, err := ot.WriteTo(os.Stdout); err != nil {
		u.CancelWithError(err)
	}
}

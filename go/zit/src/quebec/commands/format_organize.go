package commands

import (
	"flag"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type FormatOrganize struct {
	Flags organize_text.Flags
}

func init() {
	registerCommand(
		"format-organize",
		func(f *flag.FlagSet) Command {
			c := &FormatOrganize{
				Flags: organize_text.MakeFlags(),
			}

			c.Flags.AddToFlagSet(f)

			return c
		},
	)
}

func (c *FormatOrganize) Run(u *env.Env, args ...string) (err error) {
	c.Flags.Config = u.GetConfig()

	if len(args) != 1 {
		err = errors.Errorf("expected exactly one input argument")
		return
	}

	var f *os.File

	if f, err = files.Open(args[0]); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	var ot *organize_text.Text

	readOrganizeTextOp := user_ops.ReadOrganizeFile{}

	var repoId ids.RepoId

	if ot, err = readOrganizeTextOp.Run(
		u,
		f,
		repoId,
		organize_text.NewMetadata(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	ot.Options = c.Flags.GetOptionsWithMetadata(
		u.GetConfig().PrintOptions,
		u.SkuFormatBox(),
		u.GetStore().GetAbbrStore().GetAbbr(),
		u.GetExternalLikePoolForRepoId(repoId),
		ot.Metadata,
	)

	if err = ot.Refine(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = ot.WriteTo(os.Stdout); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
